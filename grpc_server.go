package main

import (
	pb "agent/grpc/service"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
)

type GrpcServerConfig struct {
	Port int `mapstructure:"port"`
}

type server struct {
	pb.UnimplementedStatusServiceServer
}

func (s *server) Heartbeat(ctx context.Context, rq *pb.Rq) (*pb.Rs, error) {
	return &pb.Rs{Status: 1}, nil
}

func StartGrpcServer(grpcServerConfig GrpcServerConfig) {
	address := "localhost:" + strconv.Itoa(grpcServerConfig.Port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		panic(fmt.Errorf("failed to listen port: %s", err))
	}
	grpcServer := grpc.NewServer()
	pb.RegisterStatusServiceServer(grpcServer, &server{})
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
