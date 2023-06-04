package main

import (
	grpc2 "agent/grpc"
	"agent/grpc/acceptor"
	pb "agent/grpc/service"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
)

const healthCheckService = "grpc.health.v1.Health"

const localhost = "localhost:"

const listenProtocol = "tcp"

type GrpcServerConfig struct {
	Port int `mapstructure:"port"`
}

type server struct {
	pb.UnimplementedBiRequestStreamServer
}

func (s *server) RequestBiStream(stream pb.BiRequestStream_RequestBiStreamServer) error {
	grpc2.Mutex.Lock()
	grpc2.RegisterConnection(stream)
	grpc2.Mutex.Unlock()
	for {
		rq, err := grpc2.GetConnection().Recv()
		if err != nil {
			return err
		}
		log.Printf("Received message: %s", rq)
		acceptor.RequestAcceptor(rq)
	}
}

func StartGrpcServer(grpcServerConfig GrpcServerConfig) {
	address := localhost + strconv.Itoa(grpcServerConfig.Port)
	listen, err := net.Listen(listenProtocol, address)
	if err != nil {
		panic(fmt.Errorf("failed to listen port: %s", err))
	}
	// 启动grpc服务
	grpcServer := grpc.NewServer()
	// 注册双向流服务
	pb.RegisterBiRequestStreamServer(grpcServer, &server{})
	// 注册健康检查服务
	healthCheckServer := health.NewServer()
	healthCheckServer.SetServingStatus(healthCheckService, healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, healthCheckServer)
	// 支持服务发现和调试
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
