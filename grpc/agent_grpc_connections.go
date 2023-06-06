package grpc

import (
	"agent/grpc/service"
	"sync"
)

var connections service.BiRequestStream_RequestBiStreamServer
var Mutex sync.Mutex

func RegisterConnection(stream service.BiRequestStream_RequestBiStreamServer) {
	if connections != nil {
		return
	}
	connections = stream
}

func GetConnection() service.BiRequestStream_RequestBiStreamServer {
	return connections
}
