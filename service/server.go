package service

import (
	"context"

	pb "github.com/ynishi/gdean/pb"
)

// server impl
type Server struct {
	pb.GDeanServiceServer
}

type ServerWithRepo struct {
	Server *Server
	Repo   ReportRepository
}

func DefaultServer() *Server {
	return &Server{}
}

func DefaultServerWithRepo(ctx context.Context, server *Server, report ReportRepository) *ServerWithRepo {
	return &ServerWithRepo{Server: server, Repo: report}
}

func (s *Server) ReportMaxEmvResults(ctx context.Context, in *pb.ReportRequest) (*pb.ReportResponse, error) {
	report := Repo.Fetch()

	return &pb.ReportResponse{Report: report}, nil
}

// helper
func NewReport() *pb.Report {
	return &pb.Report{
		NumberOfCalc: 0,
		Result:       make([]*pb.Result, 0),
	}
}
