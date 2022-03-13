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
		s.NewSqliteAnalyzeRepository, wire.Bind(new(s.AnalyzeRepository), new(*s.SqliteAnalyzeRepository)),
		s.DefaultAnalyzeServerWithRepo,
	)
	return &s.AnalyzeServer{}
}

func initializeAnalyzeServerWithMysqlRepo(ctx context.Context) *s.AnalyzeServer {
	wire.Build(
		s.DefaultMysqlAnalyzeConnInfo,
		s.NewMysqlAnalyzeRepository, wire.Bind(new(s.AnalyzeRepository), new(*s.MysqlAnalyzeRepository)),
		s.DefaultAnalyzeServerWithRepo,
	)
	return &s.AnalyzeServer{}
}
