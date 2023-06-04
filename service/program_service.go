package program_service

import (
	"agent/grpc"
	pb "agent/grpc/service"
	"agent/plugin_manager"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	"sync"
)

type PluginService struct{}

var once sync.Once

func (*PluginService) Handler(rq *pb.Payload) {

}

type Process struct {
}

func SendProgramChangeRequest(programs []plugin_manager.ProgramRs) {
	// 将字典转换为 JSON 字符串
	dataBytes, err := json.Marshal(programs)
	if err != nil {
		fmt.Println("Failed to marshal dictionary to Bytes:", err)
		return
	}

	// 创建一个 Metadata 对象
	metadata := &pb.Metadata{
		Type:     "example",
		ClientIp: "127.0.0.1",
		Headers:  map[string]string{"Header1": "Value1", "Header2": "Value2"},
	}

	// 创建一个 Any 对象
	anyData, err := anypb.New(&anypb.Any{
		Value: dataBytes,
	})

	// 创建一个 Payload 对象
	payload := &pb.Payload{
		Metadata: metadata,
		Body:     anyData,
	}

	connection := grpc.GetConnection()
	if connection == nil {
		log.Println("No connection found, no processing for this program change.")
		return
	}

	connection.Send(payload)
}

func (*PluginService) GetType() string {
	return "subscriptionProgram"
}
