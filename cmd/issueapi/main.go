//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	pb "github.com/ynishi/gdean/pb/v1"
	s "github.com/ynishi/gdean/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	suger := logger.Sugar()

	c, err := s.InitConfig()
	if err != nil {
		suger.Fatalw("failed to setup config", "err", err)
	}
	listenPort, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Port))
	if err != nil {
		suger.Fatalw("failed to listen", "err", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	asvs := initializeIssueServerWithRepo(ctx)
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
			return status.Errorf(codes.Unknown, "panic triggered: %v", p)
		}),
	}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_zap.UnaryServerInterceptor(logger),
				grpc_prometheus.UnaryServerInterceptor,
				grpc_recovery.UnaryServerInterceptor(opts...),
			),
		),
	)
	pb.RegisterIssueServiceServer(server, asvs)
	grpc_prometheus.Register(server)

	http.Handle("/metrics", promhttp.Handler())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		suger.Infow("grpc setup finished", "port", c.Port)
		server.Serve(listenPort)
	}()
	wg.Add(2)
	go func() {
		suger.Infow("metrics started", "port", 8080)
		http.ListenAndServe(":8080", nil)
	}()
	wg.Wait()
}
