package grpc

import (
	pb "agent/grpc/service"
	"log"
)

type RequestBiStreamAcceptor interface {
	Handler(rq *pb.Payload)
	GetType() string
}

var requestBiStreamAcceptors []RequestBiStreamAcceptor

func RegisterRequestBiStreamAcceptor(handler RequestBiStreamAcceptor) {
	requestBiStreamAcceptors = append(requestBiStreamAcceptors, handler)
}

func RequestAcceptor(rq *pb.Payload) {
	for _, acceptor := range requestBiStreamAcceptors {
		if acceptor.GetType() == rq.Metadata.Type {
			acceptor.Handler(rq)
		}
		log.Println("The request does not know the specific acceptor.")
	}
}
