//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/ynishi/gdean/pb"
	s "github.com/ynishi/gdean/service"
	"google.golang.org/grpc"
)

func main() {
	port := 50051
	listenPort, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	svs := initializeServerWithRepo(ctx)
	s.Repo = svs.Repo
	server := grpc.NewServer()
	pb.RegisterGDeanServiceServer(server, svs.Server)

	server.Serve(listenPort)
}
