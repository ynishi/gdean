package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/iancoleman/strcase"
	fmutil "github.com/mennanov/fieldmask-utils"
	"github.com/thanos-io/thanos/pkg/runutil"
	pb "github.com/ynishi/gdean/pb/v1"
	"google.golang.org/genproto/googleapis/rpc/code"
	gstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type IssueServiceServer struct {
	// embed service server interface, to check whether all interface impl, comment out and embed pb.UnsafeIssueServiceServer interface
	pb.IssueServiceServer
	Repo *IssueRepository
}

type UserServiceServer struct {
	pb.UserServiceServer
	// TODO: Separate from issue repository
	Repo *IssueRepository
}

func DefaultUserServiceServerWithRepo(ctx context.Context, repo *IssueRepository) *UserServiceServer {
	server := UserServiceServer{}
	if err := repo.Init(); err != nil {
		return nil
	}
	server.Repo = repo
	return &server
}

func (s *UserServiceServer) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	id, err := NewId()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if in.UserId != nil {
		id = *in.UserId
	}
	in.User.Id = &id
	user, err := s.Repo.CreateUser(ctx, in.User)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CreateUserResponse{
		User: user,
		Status: &gstatus.Status{
			Code:    int32(code.Code_OK),
			Message: "create user completed",
		},
	}, nil
}

func DefaultIssueServiceServerWithRepo(ctx context.Context, repo *IssueRepository) *IssueServiceServer {
	server := IssueServiceServer{}
	if err := repo.Init(); err != nil {
		return nil
	}
	server.Repo = repo
	return &server
}

func DefaultIssueRepository() *IssueRepository {
	return &IssueRepository{
		Rinfo: "localhost:28015",
	}
}

type IssueRepository struct {
	RSess []*r.Session
	Rinfo string
	Cache *lru.Cache
}

const ISSUEDB = "issuedb"
const ISSUETABLE = "issues"
const DATATABLE = "data"
const INTABLE = "internals"
const USERDB = "userdb"
const USERTABLE = "users"

var dbs = map[string][]string{
	ISSUEDB: {
		ISSUETABLE,
		DATATABLE,
		INTABLE,
	},
	USERDB: {
		USERTABLE,
	},
}

// Add item to repository cache with id(internal recognize its type), value as pointer.
// When query passed as id, it maybe work even no support query base cache.
func (s *IssueRepository) AddCache(id string, vp interface{}) {
	b, err := json.Marshal(vp)
	if err != nil {
		return
	}
	if s.Cache != nil {
		s.Cache.Add(fmt.Sprintf("%T-%s", vp, id), b)
	}
}

// Get item from repository cache with id(internal, recognize its type), value as pointer
func (s *IssueRepository) GetCache(id string, v interface{}) (interface{}, bool) {

	if s.Cache != nil {
		if got, ok := s.Cache.Get(fmt.Sprintf("%T-%s", v, id)); ok {
			if err := json.Unmarshal(got.([]byte), v); err != nil {
				return nil, false
			}
			return v, true
		}
	}
	return v, false
}

func (s *IssueRepository) DeleteCache(id string, v interface{}) {
	if s.Cache != nil {
		s.Cache.Remove(fmt.Sprintf("%T-%s", v, id))
	}
}

func (s *IssueRepository) Init() (err error) {
	r.SetTags("json")
	session, err := r.Connect(r.ConnectOpts{
		Address: s.Rinfo,
	})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			session.Close()
		}
	}()
	rSess := make([]*r.Session, 0)
	if len(s.RSess) > 0 {
		rSess = s.RSess
	}
	s.RSess = append(rSess, session)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	// ignore if exist, TODO: change to error check
	_ = runutil.Retry(5*time.Second, ctx.Done(), func() error {
		res, err := r.DBList().SetDifference([]interface{}{"rethinkdb", "test"}).Run(session)
		if err != nil {
			return err
		}
		var dbList []string
		if err = res.All(&dbList); err != nil {
			return err
		}
		var dbm = map[string]bool{}
		for _, v := range dbList {
			dbm[v] = true
		}
		for db, tables := range dbs {
			if _, ok := dbm[db]; !ok {
				err = r.DBCreate(db).Exec(session)
				if err != nil {
					return err
				}
			}
			res, err := r.DB(db).TableList().Run(session)
			if err != nil {
				return err
			}
			var tableList []string
			if err = res.All(&tableList); err != nil {
				return err
			}
			var tablem = map[string]bool{}
			for _, v := range tableList {
				tablem[v] = true
			}
			for _, table := range tables {
				if _, ok := tablem[table]; !ok {
					err = r.DB(db).TableCreate(table).Exec(session)
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
	cache, err := lru.New(128)
	if err != nil {
		return err
	}
	s.Cache = cache
	return nil
}

// Repository
// User

func (s *IssueRepository) CreateUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	if user.Id == nil {
		return nil, fmt.Errorf("invalid user, Id is not set")
	}
	session := s.RSess[0]
	ctxRet, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	err := runutil.Retry(2*time.Second, ctxRet.Done(), func() error {
		return r.DB(USERDB).Table(USERTABLE).Insert(user).Exec(session)
	})
	if err != nil {
		return nil, err
	}
	s.AddCache(*user.Id, user)
	return user, nil
}

func (s *IssueRepository) FetchUser(ctx context.Context, userId string) (*pb.User, error) {
	var user pb.User
	if v, ok := s.GetCache(userId, &user); ok {
		return v.(*pb.User), nil
	}
	session := s.RSess[0]
	res, err := r.DB(USERDB).Table(USERTABLE).Filter(map[string]string{
		"id": userId,
	}).Run(session)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	if err := res.One(&user); err != nil {
		return nil, err
	}
	s.AddCache(userId, &user)
	return &user, nil
}

func (s *IssueRepository) HardDeleteUser(ctx context.Context, userId string) error {
	session := s.RSess[0]
	ctxRet, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	var err error
	err = runutil.Retry(2*time.Second, ctxRet.Done(), func() error {
		_, err = r.DB(USERDB).Table(USERTABLE).Get(userId).Delete().RunWrite(session)
		return err
	})
	if err != nil {
		return err
	}
	s.DeleteCache(userId, &pb.User{})
	return nil
}

// Issue
func (s *IssueRepository) CreateIssue(ctx context.Context, issue *pb.Issue) (*pb.Issue, error) {
	session := s.RSess[0]
	ctxRet, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	err := runutil.Retry(2*time.Second, ctxRet.Done(), func() error {
		if err := r.DB(ISSUEDB).Table(ISSUETABLE).Insert(issue).Exec(session); err != nil {
			return err
		}
		return s.mutDelIssue(ctx, *issue.Id, false)
	})
	if err != nil {
		return nil, err
	}
	s.AddCache(*issue.Id, issue)
	return issue, nil
}

func (s *IssueRepository) FetchIssue(ctx context.Context, issueId string) (*pb.Issue, error) {
	var issue pb.Issue
	var user *pb.User
	if v, ok := s.GetCache(issueId, &issue); ok {
		if issue.Author == nil || issue.Author.Id == nil {
			return nil, fmt.Errorf("invalid author in cache")
		}
		user, err := s.FetchUser(ctx, *issue.Author.Id)
		if err != nil {
			return nil, err
		}
		issuep := v.(*pb.Issue)
		issuep.Author = user
		return issuep, nil
	}
	session := s.RSess[0]
	ctxRet, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	var res *r.Cursor
	var err error
	err = runutil.Retry(2*time.Second, ctxRet.Done(), func() error {
		res, err = r.DB(ISSUEDB).Table(ISSUETABLE).Filter(map[string]interface{}{
			"id":         issueId,
			"is_deleted": false,
		}).Run(session)
		return err
	})
	if err != nil {
		return nil, err
	}
	defer res.Close()
	if err := res.One(&issue); err != nil {
		return nil, err
	}
	//	return &issue, err
	if issue.Author == nil || issue.Author.Id == nil {
		return nil, fmt.Errorf("invalid author")
	}
	user, err = s.FetchUser(ctx, *issue.Author.Id)
	if err != nil {
		return nil, err
	}
	issue.Author = user
	if issue.Branches != nil || len(issue.Branches) > 0 {
		b := []*pb.Branch{}
		for i, v := range issue.Branches {
			if !v.IsDeleted {
				b = append(b, issue.Branches[i])
			}
		}
		issue.Branches = b
	}
	if issue.Comments != nil || len(issue.Comments) > 0 {
		c := []*pb.Comment{}
		for i, v := range issue.Comments {
			if !v.IsDeleted {
				c = append(c, issue.Comments[i])
			}
		}
		issue.Comments = c
	}
	s.AddCache(issueId, &issue)
	return &issue, nil
}

func (s *IssueRepository) FetchIssues(ctx context.Context, userId string) ([]*pb.Issue, error) {
	// TODO: consider use cache(maybe support query or condition)
	var user *pb.User
	var err error
	user, err = s.FetchUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	session := s.RSess[0]
	ctxRet, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	var res *r.Cursor
	err = runutil.Retry(2*time.Second, ctxRet.Done(), func() error {
		res, err = r.DB(ISSUEDB).Table(ISSUETABLE).Filter(r.Row.Field("author").Field("id").Eq(userId)).Run(session)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var issues []*pb.Issue
	if err := res.All(&issues); err != nil {
		return nil, err
	}
	for i := range issues {
		issues[i].Author = user
	}
	return issues, nil
}

func (s *IssueRepository) PutIssue(ctx context.Context, issueId string, issue *pb.Issue) (*pb.Issue, error) {
	session := s.RSess[0]
	ctxRet, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	var err error
	err = runutil.Retry(2*time.Second, ctxRet.Done(), func() error {
		_, err = r.DB(ISSUEDB).Table(ISSUETABLE).Get(issueId).Update(
			issue,
		).RunWrite(session)
		return err
	})
	if err != nil {
		return nil, err
	}
	s.AddCache(issueId, issue)
	return issue, nil
}

func (s *IssueRepository) mutDelIssue(ctx context.Context, issueId string, isDeleted bool) error {
	session := s.RSess[0]
	ctxRet, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	var err error
	err = runutil.Retry(2*time.Second, ctxRet.Done(), func() error {
		_, err = r.DB(ISSUEDB).Table(ISSUETABLE).Get(issueId).Update(map[string]interface{}{
			"is_deleted": isDeleted,
		}).RunWrite(session)
		return err
	})
	if err != nil {
		return err
	}
	s.DeleteCache(issueId, &pb.Issue{})
	return nil
}

func (s *IssueRepository) DeleteIssue(ctx context.Context, issueId string) error {
	return s.mutDelIssue(ctx, issueId, true)
}

func (s *IssueRepository) UnDeleteIssue(ctx context.Context, issueId string) error {
	return s.mutDelIssue(ctx, issueId, false)
}

func (s *IssueRepository) HardDeleteIssue(ctx context.Context, issueId string) error {
	session := s.RSess[0]
	ctxRet, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	var err error
	err = runutil.Retry(2*time.Second, ctxRet.Done(), func() error {
		_, err = r.DB(ISSUEDB).Table(ISSUETABLE).Get(issueId).Delete().RunWrite(session)
		return err
	})
	if err != nil {
		return err
	}
	s.DeleteCache(issueId, &pb.Issue{})
	return nil
}

func (s *IssueRepository) mutIssueInternal(ctx context.Context, issueId, childId, childType string, is_deleted bool) error {
	session := s.RSess[0]
	ctxRet, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	var err error
	// TODO
	err = runutil.Retry(2*time.Second, ctxRet.Done(), func() error {
		_, err = r.DB(ISSUEDB).Table(ISSUETABLE).Get(issueId).Update(map[string]interface{}{
			childType: r.Row.Field(childType).Map(func(c r.Term) interface{} {
				return r.Branch(c.Field("id").Eq(childId), c.Merge(map[string]bool{
					"is_deleted": is_deleted,
				}), c)
			}),
		}).RunWrite(session)
		return err
	})
	if err != nil {
		return err
	}
	s.DeleteCache(issueId, &pb.Issue{})
	return nil
}

func (s *IssueRepository) DeleteIssueInternal(ctx context.Context, issueId, childId, childType string) error {
	return s.mutIssueInternal(ctx, issueId, childId, childType, true)
}

func (s *IssueRepository) UnDeleteIssueInternal(ctx context.Context, issueId, childId, childType string) error {
	return s.mutIssueInternal(ctx, issueId, childId, childType, false)
}

// data for repository
func (s *IssueRepository) CreateData(ctx context.Context, data *pb.Data) (*pb.Data, error) {
	session := s.RSess[0]
	ctxRet, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	err := runutil.Retry(2*time.Second, ctxRet.Done(), func() error {
		if err := r.DB(ISSUEDB).Table(DATATABLE).Insert(data).Exec(session); err != nil {
			return err
		}
		return s.mutDelData(ctx, *data.Id, false)
	})
	if err != nil {
		return nil, err
	}
	s.AddCache(*data.Id, data)
	return data, nil
}

func (s *IssueRepository) FetchData(ctx context.Context, dataId string) (*pb.Data, error) {
	var data pb.Data
	if v, ok := s.GetCache(dataId, &data); ok {
		return v.(*pb.Data), nil
	}
	session := s.RSess[0]
	ctxRet, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	var res *r.Cursor
	var err error
	err = runutil.Retry(2*time.Second, ctxRet.Done(), func() error {
		res, err = r.DB(ISSUEDB).Table(DATATABLE).Filter(map[string]interface{}{
			"id":         dataId,
			"is_deleted": false,
		}).Run(session)
		return err
	})
	if err != nil {
		return nil, err
	}
	defer res.Close()
	if err := res.One(&data); err != nil {
		return nil, err
	}
	s.AddCache(dataId, &data)
	return &data, nil
}

func (s *IssueRepository) FetchDataList(ctx context.Context, userId string) (data []*pb.Data, err error) {
	// TODO: consider use cache(maybe support query or condition)
	session := s.RSess[0]
	ctxRet, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	var res *r.Cursor
	err = runutil.Retry(2*time.Second, ctxRet.Done(), func() error {
		res, err = r.DB(ISSUEDB).Table(DATATABLE).Filter(r.Row.Field("author").Eq(userId)).Run(session)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	defer res.Close()
	if err := res.All(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *IssueRepository) PutData(ctx context.Context, dataId string, data *pb.Data) (*pb.Data, error) {
	session := s.RSess[0]
	ctxRet, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	var err error
	err = runutil.Retry(2*time.Second, ctxRet.Done(), func() error {
		_, err = r.DB(ISSUEDB).Table(DATATABLE).Get(dataId).Update(
			data,
		).RunWrite(session)
		return err
	})
	if err != nil {
		return nil, err
	}
	s.AddCache(dataId, data)
	return data, nil
}

func (s *IssueRepository) mutDelData(ctx context.Context, dataId string, isDeleted bool) error {
	session := s.RSess[0]
	ctxRet, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	var err error
	err = runutil.Retry(2*time.Second, ctxRet.Done(), func() error {
		_, err = r.DB(ISSUEDB).Table(DATATABLE).Get(dataId).Update(map[string]interface{}{
			"is_deleted": isDeleted,
		}).RunWrite(session)
		return err
	})
	if err != nil {
		return err
	}
	s.DeleteCache(dataId, &pb.Data{})
	return nil
}

func (s *IssueRepository) DeleteData(ctx context.Context, dataId string) error {
	return s.mutDelData(ctx, dataId, true)
}

func (s *IssueRepository) UnDeleteData(ctx context.Context, dataId string) error {
	return s.mutDelData(ctx, dataId, false)
}

func (s *IssueRepository) HardDeleteData(ctx context.Context, dataId string) error {
	session := s.RSess[0]
	ctxRet, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	var err error
	err = runutil.Retry(2*time.Second, ctxRet.Done(), func() error {
		_, err = r.DB(ISSUEDB).Table(DATATABLE).Get(dataId).Delete().RunWrite(session)
		return err
	})
	if err != nil {
		return err
	}
	s.DeleteCache(dataId, &pb.User{})
	return nil
}

// service
func (s *IssueServiceServer) CreateIssue(ctx context.Context, in *pb.CreateIssueRequest) (res *pb.CreateIssueResponse, err error) {
	id, err := NewId()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// set issueId when set, no existence check, simply fail when invalid.
	if in.IssueId != nil {
		id = *in.IssueId
	}
	user, err := s.Repo.FetchUser(ctx, in.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	issue := in.Issue
	issue.Id = &id
	issue.Author = user
	issue.CreateTime = timestamppb.Now()
	resp, err := s.Repo.CreateIssue(ctx, issue)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CreateIssueResponse{
		UserId: in.UserId,
		Issue:  resp,
		Status: &gstatus.Status{
			Code:    int32(code.Code_OK),
			Message: "create issue completed",
		},
	}, nil
}

func (s *IssueServiceServer) GetIssue(ctx context.Context, in *pb.GetIssueRequest) (res *pb.GetIssueResponse, err error) {
	issue, err := s.Repo.FetchIssue(ctx, in.IssueId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.GetIssueResponse{
		Issue: issue,
		Status: &gstatus.Status{
			Code:    int32(code.Code_OK),
			Message: "get issue completed",
		},
	}, nil
}

func (s *IssueServiceServer) ListIssues(ctx context.Context, in *pb.ListIssuesRequest) (res *pb.ListIssuesResponse, err error) {
	issues, err := s.Repo.FetchIssues(ctx, in.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.ListIssuesResponse{
		Issues: issues,
		Status: &gstatus.Status{
			Code:    int32(code.Code_OK),
			Message: "list issue completed",
		},
	}, nil
}

func naming(s string) string {
	return strcase.ToCamel(s)
}

func (s *IssueServiceServer) UpdateIssue(ctx context.Context, in *pb.UpdateIssueRequest) (res *pb.UpdateIssueResponse, err error) {
	issue, err := s.Repo.FetchIssue(ctx, *in.Issue.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	mask, _ := fmutil.MaskFromProtoFieldMask(in.FieldMask, naming)
	err = fmutil.StructToStruct(mask, in.Issue, issue)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if issue, err = s.Repo.PutIssue(ctx, *in.Issue.Id, issue); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UpdateIssueResponse{
		IssueId: *in.Issue.Id,
		Issue:   issue,
		Status: &gstatus.Status{
			Code:    int32(code.Code_OK),
			Message: "update issue completed",
		},
	}, nil
}

func (s *IssueServiceServer) DeleteIssue(ctx context.Context, in *pb.DeleteIssueRequest) (res *pb.DeleteIssueResponse, err error) {
	if _, err = s.Repo.FetchIssue(ctx, in.IssueId); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	err = s.Repo.DeleteIssue(ctx, in.IssueId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteIssueResponse{
		Status: &gstatus.Status{
			Code:    int32(code.Code_OK),
			Message: "delete issue completed",
		},
	}, nil
}

func (s *IssueServiceServer) UnDeleteIssue(ctx context.Context, in *pb.UnDeleteIssueRequest) (res *pb.UnDeleteIssueResponse, err error) {
	err = s.Repo.UnDeleteIssue(ctx, in.IssueId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UnDeleteIssueResponse{
		Status: &gstatus.Status{
			Code:    int32(code.Code_OK),
			Message: "delete issue completed",
		},
	}, nil
}

func (s *IssueServiceServer) ExportIssue(ctx context.Context, in *pb.ExportIssueRequest) (res *pb.ExportIssueResponse, err error) {
	issues, err := s.Repo.FetchIssues(ctx, in.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.ExportIssueResponse{
		Issues: issues,
		Status: &gstatus.Status{
			Code:    int32(code.Code_OK),
			Message: "list issue completed",
		},
	}, nil
}

func (s *IssueServiceServer) DeleteIssueInternal(ctx context.Context, in *pb.DeleteIssueInternalRequest) (res *pb.DeleteIssueInternalResponse, err error) {
	if _, err = s.Repo.FetchIssue(ctx, in.IssueId); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err = s.Repo.DeleteIssueInternal(ctx, in.IssueId, in.ChildId, in.ChildType); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteIssueInternalResponse{
		IssueId: in.IssueId,
		Status: &gstatus.Status{
			Code:    int32(code.Code_OK),
			Message: "delete issue internal completed",
		},
	}, nil
}

func (s *IssueServiceServer) UnDeleteIssueInternal(ctx context.Context, in *pb.UnDeleteIssueInternalRequest) (res *pb.UnDeleteIssueInternalResponse, err error) {
	if _, err = s.Repo.FetchIssue(ctx, in.IssueId); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err = s.Repo.UnDeleteIssueInternal(ctx, in.IssueId, in.ChildId, in.ChildType); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UnDeleteIssueInternalResponse{
		IssueId: in.IssueId,
		Status: &gstatus.Status{
			Code:    int32(code.Code_OK),
			Message: "undelete issue internal completed",
		},
	}, nil
}

// data
func (s *IssueServiceServer) CreateData(ctx context.Context, in *pb.CreateDataRequest) (*pb.CreateDataResponse, error) {
	id, err := NewId()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// set issueId when set, no existence check, simply fail when invalid.
	if in.DataId != nil {
		id = *in.DataId
	}
	if _, err := s.Repo.FetchUser(ctx, in.UserId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	data := in.Data
	data.Id = &id
	data.CreateTime = timestamppb.Now()

	if _, err = s.Repo.CreateData(ctx, data); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CreateDataResponse{
		Data: data,
		Status: &gstatus.Status{
			Code:    int32(code.Code_OK),
			Message: "create data completed",
		},
	}, nil
}

func (s *IssueServiceServer) GetData(ctx context.Context, in *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	data, err := s.Repo.FetchData(ctx, in.DataId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.GetDataResponse{
		Data: data,
		Status: &gstatus.Status{
			Code:    int32(code.Code_OK),
			Message: "get data completed",
		},
	}, nil
}

func (s *IssueServiceServer) ListData(ctx context.Context, in *pb.ListDataRequest) (*pb.ListDataResponse, error) {
	data, err := s.Repo.FetchDataList(ctx, in.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.ListDataResponse{
		Data: data,
		Status: &gstatus.Status{
			Code:    int32(code.Code_OK),
			Message: "list data completed",
		},
	}, nil
}

func (s *IssueServiceServer) UpdateData(ctx context.Context, in *pb.UpdateDataRequest) (*pb.UpdateDataResponse, error) {
	if _, err := s.Repo.FetchData(ctx, *in.Data.Id); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	data, err := s.Repo.PutData(ctx, *in.Data.Id, in.Data)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UpdateDataResponse{
		Data: data,
		Status: &gstatus.Status{
			Code:    int32(code.Code_OK),
			Message: "list data completed",
		},
	}, nil
}

func (s *IssueServiceServer) DeleteData(ctx context.Context, in *pb.DeleteDataRequest) (*pb.DeleteDataResponse, error) {
	if _, err := s.Repo.FetchData(ctx, in.DataId); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err := s.Repo.DeleteData(ctx, in.DataId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteDataResponse{
		Status: &gstatus.Status{
			Code:    int32(code.Code_OK),
			Message: "delete data completed",
		},
	}, nil
}

func (s *IssueServiceServer) UnDeleteData(ctx context.Context, in *pb.UnDeleteDataRequest) (*pb.UnDeleteDataResponse, error) {
	if _, err := s.Repo.FetchData(ctx, in.DataId); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err := s.Repo.UnDeleteData(ctx, in.DataId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UnDeleteDataResponse{
		Status: &gstatus.Status{
			Code:    int32(code.Code_OK),
			Message: "delete data completed",
		},
	}, nil
}

func (s *IssueServiceServer) DecideBranch(ctx context.Context, in *pb.DecideBranchRequest) (*pb.DecideBranchResponse, error) {

	issue, err := s.Repo.FetchIssue(ctx, in.IssueId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	issue.Decided = in.BranchId
	if _, err := s.Repo.PutIssue(ctx, in.IssueId, issue); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DecideBranchResponse{
		Status: &gstatus.Status{
			Code:    int32(code.Code_OK),
			Message: "decide branch completed",
		},
	}, nil
}

func (s *IssueServiceServer) AddAnalyzedResult(ctx context.Context, in *pb.AddAnalyzedResultRequest) (*pb.AddAnalyzedResultResponse, error) {
	issue, err := s.Repo.FetchIssue(ctx, in.IssueId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	issue.Results = append(issue.Results, in.AnalyzedResult)
	if _, err := s.Repo.PutIssue(ctx, in.IssueId, issue); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.AddAnalyzedResultResponse{
		Status: &gstatus.Status{
			Code:    int32(code.Code_OK),
			Message: "add analyze result completed",
		},
	}, nil
}
