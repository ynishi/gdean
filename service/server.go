package service

import (
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
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

func (s *Server) MaxEmv(ctx context.Context, in *pb.MaxEmvRequest) (*pb.MaxEmvResponse, error) {
	log.Printf("Recieved: %v", in.GetTowPData().GetP1())
	maxEmv, _ := calcMaxEmv(in.GetTowPData().GetP1(), in.GetTowPData().GetDataP1(), in.GetTowPData().GetDataP2())
	now := time.Now()
	result := pb.Result{MaxEmv: maxEmv, CreateTime: &timestamp.Timestamp{
		Seconds: now.Unix(),
		Nanos:   int32(now.Nanosecond()),
	}}
	Repo.Put(&result)
	return &pb.MaxEmvResponse{Result: &result}, nil
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
