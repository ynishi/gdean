package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	pb "github.com/ynishi/gdean/pb/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

var testAddress = "localhost:28115"

func NewMockIssueRepository() *IssueRepository {
	repo := IssueRepository{Rinfo: testAddress}
	err := repo.Init()
	if err != nil {
		panic(err)
	}
	return &repo
}

func TestShouldCreateAndFetchIssue(t *testing.T) {
	repo := NewMockIssueRepository()
	// don't use cache
	repo.Cache = nil
	assert.NotNil(t, repo.RSess)

	id, err := NewId()
	assert.Nil(t, err)

	uid, err := NewId()
	assert.Nil(t, err)

	author := pb.User{Id: &uid, Name: "test", CreateTime: timestamppb.Now()}
	issue := pb.Issue{Id: &id, Title: "test_title", Desc: "test_desc", Author: &author, CreateTime: timestamppb.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// clean up Issue
	_ = repo.HardDeleteIssue(ctx, *issue.Id)

	uRet, err := repo.CreateUser(ctx, &author)
	assert.Nil(t, err)
	assert.NotNil(t, uRet)
	assert.Equal(t, author.Id, uRet.Id)

	ret, err := repo.CreateIssue(ctx, &issue)
	assert.Nil(t, err)
	assert.NotNil(t, ret)
	assert.Equal(t, &issue, ret)

	got, err := repo.FetchIssue(ctx, *issue.Id)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, &issue, got)
	assert.Equal(t, issue.Id, got.Id)
	assert.Equal(t, author.Id, got.Author.Id)

	err = repo.HardDeleteIssue(ctx, *issue.Id)
	assert.Nil(t, err)
	err = repo.HardDeleteUser(ctx, *author.Id)
	assert.Nil(t, err)

}

func TestShouldDeleteInternalIssue(t *testing.T) {
	repo := NewMockIssueRepository()
	// don't use cache
	repo.Cache = nil
	assert.NotNil(t, repo.RSess)

	id, err := NewId()
	assert.Nil(t, err)

	uid, err := NewId()
	assert.Nil(t, err)

	bid1, err := NewId()
	assert.Nil(t, err)
	bid2, err := NewId()
	assert.Nil(t, err)

	author := pb.User{Id: &uid, Name: "test", CreateTime: timestamppb.Now()}
	branch1 := pb.Branch{Id: &bid1, Content: "test_branch1"}
	branch2 := pb.Branch{Id: &bid2, Content: "test_branch2"}
	branches := []*pb.Branch{&branch1, &branch2}
	issue := pb.Issue{Id: &id, Title: "test_title", Desc: "test_desc", Author: &author, Branches: branches, CreateTime: timestamppb.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// clean up Issue
	_ = repo.HardDeleteIssue(ctx, *issue.Id)

	uRet, err := repo.CreateUser(ctx, &author)
	assert.Nil(t, err)
	assert.NotNil(t, uRet)
	assert.Equal(t, author.Id, uRet.Id)

	ret, err := repo.CreateIssue(ctx, &issue)
	assert.Nil(t, err)
	assert.NotNil(t, ret)
	assert.Equal(t, &issue, ret)

	err = repo.DeleteIssueInternal(ctx, *issue.Id, *branch1.Id, "branches")
	assert.Nil(t, err)

	got, err := repo.FetchIssue(ctx, *issue.Id)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, 1, len(got.Branches))
	assert.Equal(t, issue.Branches[1], got.Branches[0])
	assert.Equal(t, issue.Id, got.Id)
	assert.Equal(t, author.Id, got.Author.Id)

	err = repo.HardDeleteIssue(ctx, *issue.Id)
	assert.Nil(t, err)
	err = repo.HardDeleteUser(ctx, *author.Id)
	assert.Nil(t, err)

}

func TestShouldUnDeleteInternalIssue(t *testing.T) {
	repo := NewMockIssueRepository()
	// don't use cache
	repo.Cache = nil
	assert.NotNil(t, repo.RSess)

	id, err := NewId()
	assert.Nil(t, err)

	uid, err := NewId()
	assert.Nil(t, err)

	bid1, err := NewId()
	assert.Nil(t, err)
	bid2, err := NewId()
	assert.Nil(t, err)

	author := pb.User{Id: &uid, Name: "test", CreateTime: timestamppb.Now()}
	branch1 := pb.Branch{Id: &bid1, Content: "test_branch1"}
	branch2 := pb.Branch{Id: &bid2, Content: "test_branch2"}
	branches := []*pb.Branch{&branch1, &branch2}
	issue := pb.Issue{Id: &id, Title: "test_title", Desc: "test_desc", Author: &author, Branches: branches, CreateTime: timestamppb.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// clean up Issue
	_ = repo.HardDeleteIssue(ctx, *issue.Id)

	uRet, err := repo.CreateUser(ctx, &author)
	assert.Nil(t, err)
	assert.NotNil(t, uRet)
	assert.Equal(t, author.Id, uRet.Id)

	ret, err := repo.CreateIssue(ctx, &issue)
	assert.Nil(t, err)
	assert.NotNil(t, ret)
	assert.Equal(t, &issue, ret)

	err = repo.DeleteIssueInternal(ctx, *issue.Id, *branch1.Id, "branches")
	assert.Nil(t, err)

	got, err := repo.FetchIssue(ctx, *issue.Id)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, 1, len(got.Branches))
	assert.Equal(t, issue.Branches[1], got.Branches[0])
	assert.Equal(t, issue.Id, got.Id)
	assert.Equal(t, author.Id, got.Author.Id)

	err = repo.UnDeleteIssueInternal(ctx, *issue.Id, *branch1.Id, "branches")
	assert.Nil(t, err)

	got2, err := repo.FetchIssue(ctx, *issue.Id)
	assert.Nil(t, err)
	assert.NotNil(t, got2)
	assert.Equal(t, 2, len(got2.Branches))
	assert.Equal(t, issue.Branches[0], got2.Branches[0])
	assert.Equal(t, issue.Id, got2.Id)
	assert.Equal(t, author.Id, got2.Author.Id)

	err = repo.HardDeleteIssue(ctx, *issue.Id)
	assert.Nil(t, err)
	err = repo.HardDeleteUser(ctx, *author.Id)
	assert.Nil(t, err)

}

func TestShouldCreateAndFetchIssueWithCache(t *testing.T) {
	repo := NewMockIssueRepository()
	assert.NotNil(t, repo.RSess)
	assert.NotNil(t, repo.Cache)

	id, err := NewId()
	assert.Nil(t, err)

	uid, err := NewId()
	assert.Nil(t, err)

	author := pb.User{Id: &uid, Name: "test1", CreateTime: timestamppb.Now()}
	issue := pb.Issue{Id: &id, Title: "test1_title", Desc: "test1_desc", Author: &author, CreateTime: timestamppb.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// clean up Issue
	_ = repo.HardDeleteIssue(ctx, *issue.Id)

	uRet, err := repo.CreateUser(ctx, &author)
	assert.Nil(t, err)
	assert.NotNil(t, uRet)
	assert.Equal(t, author.Id, uRet.Id)

	ret, err := repo.CreateIssue(ctx, &issue)
	assert.Nil(t, err)
	assert.NotNil(t, ret)
	assert.Equal(t, &issue, ret)

	got, err := repo.FetchIssue(ctx, *issue.Id)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, &issue, got)
	assert.Equal(t, issue.Id, got.Id)
	assert.Equal(t, author.Id, got.Author.Id)

	err = repo.HardDeleteIssue(ctx, *issue.Id)
	assert.Nil(t, err)
	err = repo.HardDeleteUser(ctx, *author.Id)
	assert.Nil(t, err)
}

func TestShouldFetchIssues(t *testing.T) {
	repo := NewMockIssueRepository()

	id1, _ := NewId()
	id2, _ := NewId()
	uid, _ := NewId()

	author := pb.User{Id: &uid, Name: "test1", CreateTime: timestamppb.Now()}
	issue1 := pb.Issue{Id: &id1, Title: "test1_title1", Desc: "test1_desc1", Author: &author, CreateTime: timestamppb.Now()}
	issue2 := pb.Issue{Id: &id2, Title: "test1_title2", Desc: "test1_desc2", Author: &author, CreateTime: timestamppb.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// clean up Issue
	_ = repo.HardDeleteIssue(ctx, *issue1.Id)
	_ = repo.HardDeleteIssue(ctx, *issue2.Id)

	_, err := repo.CreateUser(ctx, &author)
	assert.Nil(t, err)

	_, err = repo.CreateIssue(ctx, &issue1)
	assert.Nil(t, err)
	_, err = repo.CreateIssue(ctx, &issue2)
	assert.Nil(t, err)

	got, err := repo.FetchIssues(ctx, *author.Id)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(got))
	// allow no order
	var acc = map[string]*pb.Issue{}
	for _, v := range got {
		acc[*v.Id] = v
	}
	assert.Equal(t, &issue1, acc[*issue1.Id])
	assert.Equal(t, &issue2, acc[*issue2.Id])
	assert.Equal(t, issue1.Id, acc[*issue1.Id].Id)
	assert.Equal(t, author.Id, acc[*issue1.Id].Author.Id)

	err = repo.HardDeleteIssue(ctx, *issue1.Id)
	assert.Nil(t, err)
	err = repo.HardDeleteIssue(ctx, *issue2.Id)
	assert.Nil(t, err)
	err = repo.HardDeleteUser(ctx, *author.Id)
	assert.Nil(t, err)
}

func TestShouldPutIssue(t *testing.T) {
	repo := NewMockIssueRepository()

	id, _ := NewId()
	uid, _ := NewId()

	author := pb.User{Id: &uid, Name: "test1", CreateTime: timestamppb.Now()}
	issue1 := pb.Issue{Id: &id, Title: "test1_title1", Desc: "test1_desc1", Author: &author, CreateTime: timestamppb.Now()}
	issue2 := pb.Issue{Id: &id, Title: "test1_title2", Desc: "test1_desc2", Author: &author, CreateTime: timestamppb.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// clean up Issue
	_ = repo.HardDeleteIssue(ctx, *issue1.Id)

	_, err := repo.CreateUser(ctx, &author)
	assert.Nil(t, err)

	_, err = repo.CreateIssue(ctx, &issue1)
	assert.Nil(t, err)
	res, err := repo.PutIssue(ctx, *issue2.Id, &issue2)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, issue2.Title, res.Title)

	got, err := repo.FetchIssue(ctx, *issue2.Id)
	assert.Nil(t, err)
	assert.Equal(t, &issue2, got)
	assert.Equal(t, issue2.Id, got.Id)

	err = repo.HardDeleteIssue(ctx, id)
	assert.Nil(t, err)
	err = repo.HardDeleteUser(ctx, *author.Id)
	assert.Nil(t, err)
}

func TestShouldDeleteIssue(t *testing.T) {
	repo := NewMockIssueRepository()

	id, _ := NewId()
	uid, _ := NewId()

	author := pb.User{Id: &uid, Name: "test1", CreateTime: timestamppb.Now()}
	issue := pb.Issue{Id: &id, Title: "test1_title1", Desc: "test1_desc1", Author: &author, CreateTime: timestamppb.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// clean up Issue
	_ = repo.HardDeleteIssue(ctx, *issue.Id)

	_, err := repo.CreateUser(ctx, &author)
	assert.Nil(t, err)

	_, err = repo.CreateIssue(ctx, &issue)
	assert.Nil(t, err)

	pre, err := repo.FetchIssue(ctx, *issue.Id)
	assert.Nil(t, err)
	assert.NotNil(t, pre)

	err = repo.DeleteIssue(ctx, *issue.Id)
	assert.Nil(t, err)
	got, err := repo.FetchIssue(ctx, *issue.Id)
	assert.NotNil(t, err)
	assert.Nil(t, got)

	err = repo.HardDeleteIssue(ctx, id)
	assert.Nil(t, err)
	err = repo.HardDeleteUser(ctx, *author.Id)
	assert.Nil(t, err)
}

func TestShouldUnDeleteIssue(t *testing.T) {
	repo := NewMockIssueRepository()

	id, _ := NewId()
	uid, _ := NewId()

	author := pb.User{Id: &uid, Name: "test1", CreateTime: timestamppb.Now()}
	issue := pb.Issue{Id: &id, Title: "test1_title1", Desc: "test1_desc1", Author: &author, CreateTime: timestamppb.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// clean up Issue
	_ = repo.HardDeleteIssue(ctx, *issue.Id)

	_, err := repo.CreateUser(ctx, &author)
	assert.Nil(t, err)

	_, err = repo.CreateIssue(ctx, &issue)
	assert.Nil(t, err)

	err = repo.DeleteIssue(ctx, *issue.Id)
	assert.Nil(t, err)
	_, err = repo.FetchIssue(ctx, *issue.Id)
	assert.NotNil(t, err)

	err = repo.UnDeleteIssue(ctx, *issue.Id)
	assert.Nil(t, err)

	// expect un deleted issue available
	got, err := repo.FetchIssue(ctx, *issue.Id)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, &issue, got)
	assert.Equal(t, issue.Id, got.Id)

	err = repo.HardDeleteIssue(ctx, id)
	assert.Nil(t, err)
	err = repo.HardDeleteUser(ctx, *author.Id)
	assert.Nil(t, err)
}

func TestShouldCreateAndFetchData(t *testing.T) {
	repo := NewMockIssueRepository()
	// don't use cache
	repo.Cache = nil
	assert.NotNil(t, repo.RSess)

	id, err := NewId()
	assert.Nil(t, err)

	data := pb.Data{Id: &id, Columns: []string{"c1", "c2"}, CreateTime: timestamppb.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// clean up Issue
	_ = repo.HardDeleteData(ctx, id)

	ret, err := repo.CreateData(ctx, &data)
	assert.Nil(t, err)
	assert.NotNil(t, ret)
	assert.Equal(t, &data, ret)

	got, err := repo.FetchData(ctx, *data.Id)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, &data, got)
	assert.Equal(t, data.Id, got.Id)

	err = repo.HardDeleteData(ctx, *data.Id)
	assert.Nil(t, err)

}

func TestShouldCreateAndFetchDataWithCache(t *testing.T) {
	repo := NewMockIssueRepository()
	assert.NotNil(t, repo.RSess)
	assert.NotNil(t, repo.Cache)

	id, err := NewId()
	assert.Nil(t, err)

	data := pb.Data{Id: &id, Columns: []string{"c1", "c2"}, CreateTime: timestamppb.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// clean up Issue
	_ = repo.HardDeleteData(ctx, *data.Id)

	ret, err := repo.CreateData(ctx, &data)
	assert.Nil(t, err)
	assert.NotNil(t, ret)
	assert.Equal(t, &data, ret)

	got, err := repo.FetchData(ctx, *data.Id)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, &data, got)
	assert.Equal(t, data.Id, got.Id)

	err = repo.HardDeleteData(ctx, *data.Id)
	assert.Nil(t, err)
}

func TestShouldFetchDataList(t *testing.T) {
	repo := NewMockIssueRepository()

	id1, _ := NewId()
	id2, _ := NewId()
	uid, _ := NewId()

	data1 := pb.Data{Id: &id1, Columns: []string{"c1", "c2"}, Author: uid, CreateTime: timestamppb.Now()}
	data2 := pb.Data{Id: &id2, Columns: []string{"d1", "d2"}, Author: uid, CreateTime: timestamppb.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// clean up Issue
	_ = repo.HardDeleteData(ctx, *data1.Id)
	_ = repo.HardDeleteData(ctx, *data2.Id)

	var err error
	_, err = repo.CreateData(ctx, &data1)
	assert.Nil(t, err)
	_, err = repo.CreateData(ctx, &data2)
	assert.Nil(t, err)

	got, err := repo.FetchDataList(ctx, uid)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(got))
	// allow no order
	var acc = map[string]*pb.Data{}
	for _, v := range got {
		acc[*v.Id] = v
	}
	assert.Equal(t, &data1, acc[*data1.Id])
	assert.Equal(t, &data2, acc[*data2.Id])
	assert.Equal(t, data1.Id, acc[*data1.Id].Id)

	err = repo.HardDeleteData(ctx, *data1.Id)
	assert.Nil(t, err)
	err = repo.HardDeleteData(ctx, *data2.Id)
	assert.Nil(t, err)
}

func TestShouldPutData(t *testing.T) {
	repo := NewMockIssueRepository()

	id, _ := NewId()
	uid, _ := NewId()

	data1 := pb.Data{Id: &id, Columns: []string{"c1", "c2"}, Author: uid, CreateTime: timestamppb.Now()}
	data2 := pb.Data{Id: &id, Columns: []string{"d1", "d2"}, Author: uid, CreateTime: timestamppb.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// clean up Data
	_ = repo.HardDeleteData(ctx, *data1.Id)

	_, err := repo.CreateData(ctx, &data1)
	assert.Nil(t, err)
	res, err := repo.PutData(ctx, *data2.Id, &data2)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, data2.Columns, res.Columns)

	got, err := repo.FetchData(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, &data2, got)
	assert.Equal(t, data2.Id, got.Id)

	err = repo.HardDeleteData(ctx, id)
	assert.Nil(t, err)
}

func TestShouldDeleteData(t *testing.T) {
	repo := NewMockIssueRepository()

	id, _ := NewId()
	uid, _ := NewId()

	data := pb.Data{Id: &id, Columns: []string{"c1", "c2"}, Author: uid, CreateTime: timestamppb.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// clean up Data
	_ = repo.HardDeleteData(ctx, id)

	_, err := repo.CreateData(ctx, &data)
	assert.Nil(t, err)

	pre, err := repo.FetchData(ctx, id)
	assert.Nil(t, err)
	assert.NotNil(t, pre)

	err = repo.DeleteData(ctx, id)
	assert.Nil(t, err)
	got, err := repo.FetchData(ctx, id)
	assert.NotNil(t, err)
	assert.Nil(t, got)

	err = repo.HardDeleteData(ctx, id)
	assert.Nil(t, err)
}

func TestShouldUnDeleteData(t *testing.T) {
	repo := NewMockIssueRepository()

	id, _ := NewId()
	uid, _ := NewId()

	data := pb.Data{Id: &id, Columns: []string{"c1", "c2"}, Author: uid, CreateTime: timestamppb.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// clean up Data
	_ = repo.HardDeleteData(ctx, id)

	_, err := repo.CreateData(ctx, &data)
	assert.Nil(t, err)

	err = repo.DeleteData(ctx, id)
	assert.Nil(t, err)
	_, err = repo.FetchData(ctx, id)
	assert.NotNil(t, err)

	err = repo.UnDeleteData(ctx, id)
	assert.Nil(t, err)

	// expect un deleted issue available
	got, err := repo.FetchData(ctx, id)

	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, &data, got)
	assert.Equal(t, data.Id, got.Id)

	err = repo.HardDeleteData(ctx, id)
	assert.Nil(t, err)
}
