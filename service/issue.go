package service

import (
	"context"
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
	pb.IssueServiceServer
	Repo *IssueRepository
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

var dbs = map[string][]string{
	"issuedb": {
		"issues",
	},
	"userdb": {
		"users",
	},
}

// Add item to repository cache with id(internal recognize its type), value as pointer.
// When query passed as id, it maybe work even no support query base cache.
func (s *IssueRepository) AddCache(id string, vp interface{}) {
	if s.Cache != nil {
		s.Cache.Add(fmt.Sprintf("%T-%s", vp, id), vp)
	}
}

// Get item from repository cache with id(internal, recognize its type), value as pointer
func (s *IssueRepository) GetCache(id string, v interface{}) (interface{}, bool) {
	if s.Cache != nil {
		return s.Cache.Get(fmt.Sprintf("%T-%s", v, id))
	}
	return v, false
}

func (s *IssueRepository) DeleteCache(id string, v interface{}) {
	if s.Cache != nil {
		s.Cache.Remove(fmt.Sprintf("%T-%s", v, id))
	}
}

func (s *IssueRepository) Init() (err error) {
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
	s.RSess = append(s.RSess, session)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// ignore if exist, TODO: change to error check
	_ = runutil.Retry(5*time.Second, ctx.Done(), func() error {
		for db, tables := range dbs {
			err = r.DBCreate(db).Exec(session)
			if err != nil {
				return err
			}
			for _, table := range tables {
				err = r.DB(db).TableCreate(table).Exec(session)
				if err != nil {
					return err
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
func (s *IssueRepository) FetchUser(ctx context.Context, userId string) (*pb.User, error) {
	var user pb.User
	if v, ok := s.GetCache(userId, &user); ok {
		return v.(*pb.User), nil
	}
	session := s.RSess[0]
	res, err := r.DB("userdb").Table("users").Filter(map[string]string{
		"user_id": userId,
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

// Issue
func (s *IssueRepository) CreateIssue(ctx context.Context, issue *pb.Issue) (*pb.Issue, error) {
	session := s.RSess[0]
	ctxRet, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	err := runutil.Retry(2*time.Second, ctxRet.Done(), func() error {
		return r.DB("issuedb").Table("issues").Insert(issue).Exec(session)
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
		res, err = r.DB("issuedb").Table("issues").Filter(map[string]string{
			"issue_id": issueId,
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
	user, err = s.FetchUser(ctx, *issue.Author.Id)
	if err != nil {
		return nil, err
	}
	issue.Author = user
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
		res, err = r.DB("issuedb").Table("issues").Filter(map[string]string{
			"author": userId,
		}).Run(session)
		return err
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
		_, err = r.DB("issuedb").Table("issues").Get(issueId).Update(
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

func (s *IssueRepository) DeleteIssue(ctx context.Context, issueId string) error {
	session := s.RSess[0]
	ctxRet, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	var err error
	err = runutil.Retry(2*time.Second, ctxRet.Done(), func() error {
		_, err = r.DB("issuedb").Table("issues").Get(issueId).Update(map[string]interface{}{
			"is_deleted": true,
		}).RunWrite(session)
		return err
	})
	if err != nil {
		return err
	}
	s.DeleteCache(issueId, &pb.Issue{})
	return nil
}

func (s *IssueRepository) UnDeleteIssue(ctx context.Context, issueId string) error {
	session := s.RSess[0]
	ctxRet, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	var err error
	err = runutil.Retry(2*time.Second, ctxRet.Done(), func() error {
		_, err = r.DB("issuedb").Table("issues").Get(issueId).Update(map[string]interface{}{
			"is_deleted": false,
		}).RunWrite(session)
		return err
	})
	if err != nil {
		return err
	}
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
	if _, err = s.Repo.FetchIssue(ctx, in.IssueId); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
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
