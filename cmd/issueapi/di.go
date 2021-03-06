//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	s "github.com/ynishi/gdean/service"
)

func initializeIssueServerWithRepo(ctx context.Context) *s.IssueServiceServer {
	wire.Build(
		s.DefaultIssueRepository,
		s.DefaultIssueServiceServerWithRepo,
	)
	return &s.IssueServiceServer{}
}

func initializeUserServerWithRepo(ctx context.Context) *s.UserServiceServer {
	wire.Build(
		s.DefaultIssueRepository,
		s.DefaultUserServiceServerWithRepo,
	)
	return &s.UserServiceServer{}
}
