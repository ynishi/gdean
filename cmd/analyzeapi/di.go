//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	s "github.com/ynishi/gdean/service"
)

func initializeServerWithRepo(ctx context.Context) *s.ServerWithRepo {
	wire.Build(
		s.DefaultSqlite3ConnInfo,
		s.NewSqlite3ReportRepository, wire.Bind(new(s.ReportRepository), new(*s.Sqlite3ReportRepository)),
		s.DefaultServer,
		s.DefaultServerWithRepo,
	)
	return &s.ServerWithRepo{}
}

func initializeAnalyzeServerWithRepo(ctx context.Context) *s.AnalyzeServer {
	wire.Build(
		s.DefaultSqliteAnalyzeConnInfo,
		s.DefaultAnalyzeServiceServer,
		s.AnalyzeServerWithSqliteRepo,
	)
	return &s.AnalyzeServer{}
}
