package service

import (
	"context"
	pb "github.com/ynishi/gdean/pb"
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
	return &IssueRepository{}
}

type IssueRepository struct{}

func (s *IssueRepository) Init() (err error) {
	return nil
}
