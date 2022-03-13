package service

import (
	"context"
	"fmt"
	pb "github.com/ynishi/gdean/pb"
	"google.golang.org/grpc"
	"time"
)

var DialInfo = "localhost:50051"

func MaxEmv(p1 float32, d1 []int32, d2 []int32) int32 {
	conn, _ := grpc.Dial(DialInfo, grpc.WithInsecure())
	defer conn.Close()
	c := pb.NewAnalyzeServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := pb.MaxEmvRequest{
		TowPData: &pb.TowPData{
			P1:     p1,
			DataP1: d1,
			DataP2: d2,
		},
	}
	emv, err := c.MaxEmv(ctx, &req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("got emv:%d", emv.GetMaxEmv())
	return emv.GetMaxEmv()
}

func ReportMaxEmvResults() *pb.Report {
	conn, _ := grpc.Dial(DialInfo, grpc.WithInsecure())
	defer conn.Close()
	c := pb.NewGDeanServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, _ := c.ReportMaxEmvResults(ctx, &pb.ReportRequest{})
	return r.GetReport()
}

func CreateMeta(name, desc string, isAvailable bool, paramDef map[string]string) *pb.CreateMetaResponse {
	conn, _ := grpc.Dial(DialInfo, grpc.WithInsecure())
	defer conn.Close()
	c := pb.NewAnalyzeServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := pb.CreateMetaRequest{}
	req.MetaBody = &pb.MetaBody{Name: name, Desc: desc, IsAvailable: isAvailable, ParamDef: paramDef}
	r, _ := c.CreateMeta(ctx, &req)
	return r
}

func GetMeta(id uint32) *pb.GetMetaResponse {
	conn, _ := grpc.Dial(DialInfo, grpc.WithInsecure())
	defer conn.Close()
	c := pb.NewAnalyzeServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := pb.GetMetaRequest{}
	req.Id = id
	r, _ := c.GetMeta(ctx, &req)
	return r
}
