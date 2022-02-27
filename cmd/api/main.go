//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	pb "github.com/ynishi/gdean/pb"
	s "github.com/ynishi/gdean/service"
	"google.golang.org/grpc"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	c, err := s.InitConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to setup config")
	}
	listenPort, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	svs := initializeServerWithRepo(ctx)
	s.Repo = svs.Repo
	server := grpc.NewServer()
	pb.RegisterGDeanServiceServer(server, svs.Server)
	log.Info().Int("port", c.Port).Msg("setup finished")

	server.Serve(listenPort)
}
