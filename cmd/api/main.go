package main

import (
	"fmt"
	pb "github.com/ynishi/gdean/pb"
	"github.com/ynishi/gdean/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	port := 50051
	listenPort, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	service.Repo = service.NewSqlite3ReportRepository(nil)
	server := grpc.NewServer()
	pb.RegisterGDeanServiceServer(server, &service.Server{})

	server.Serve(listenPort)
}
