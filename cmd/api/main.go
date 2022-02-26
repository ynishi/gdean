package main

import (
	"context"
	"fmt"
	pb "github.com/ynishi/gdean/pb"
	"github.com/ynishi/gdean/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	port := 50051
	listenPort, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	svs := service.InitializeServerWithRepo(ctx)
	service.Repo = svs.Repo
	server := grpc.NewServer()
	pb.RegisterGDeanServiceServer(server, svs.Server)

	server.Serve(listenPort)
}
