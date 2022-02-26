package service

import (
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/thoas/go-funk"
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

// max emv logic
func calcMaxEmv(p1 float32, dataP1 []int32, dataP2 []int32) (int32, error) {
	p2 := 1 - p1
	l := len(dataP1)
	dataP12 := make([]float32, l)
	dataP22 := make([]float32, l)
	sums := make([]float32, l)
	for i := 0; i < l; i++ {
		dataP12[i] = float32(dataP1[i]) * p1
		dataP22[i] = float32(dataP2[i]) * p2
		sums[i] = dataP12[i] + dataP22[i]
	}
	maxSum := funk.MaxFloat32(sums)
	maxSumI := funk.IndexOf(sums, maxSum)

	return int32(maxSumI), nil

}

// helper
func NewReport() *pb.Report {
	return &pb.Report{
		NumberOfCalc: 0,
		Result:       make([]*pb.Result, 0),
	}
}
