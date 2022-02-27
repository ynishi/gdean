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
	c := pb.NewGDeanServiceClient(conn)

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
	fmt.Println("got emv:" + emv.GetResult().String())
	return emv.GetResult().GetMaxEmv()
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
