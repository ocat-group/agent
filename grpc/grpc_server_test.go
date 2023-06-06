package grpc

import (
	pb "agent/grpc/service"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/anypb"
	"io"
	"log"
	"testing"
)

func TestGrpc(t *testing.T) {
	//建立无认证的连接
	conn, _ := grpc.Dial(":9001", grpc.WithTransportCredentials(insecure.NewCredentials()))

	defer conn.Close()
	client := pb.NewBiRequestStreamClient(conn)

	log.Printf("start of stream")
	requestBiStreamClient, _ := client.RequestBiStream(context.Background())
	rq, _ := anypb.New(&anypb.Any{
		Value: []byte("请求一下你"),
	})
	metadata := pb.Metadata{Type: "beta", ClientIp: "localhost"}
	_ = requestBiStreamClient.Send(&pb.Payload{Metadata: &metadata, Body: rq})
	for {
		rs, err := requestBiStreamClient.Recv()
		if err == io.EOF {
			log.Println("end of stream")
			break
		}
		log.Printf("Receive the response from the server：%s", rs)
	}
}
