package program_service

import (
	"agent/grpc"
	pb "agent/grpc/service"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	"sync"
	"time"
)

type PluginService struct{}

var once sync.Once

func init() {
	grpc.RegisterRequestBiStreamAcceptor(&PluginService{})
}

func (*PluginService) Handler(rq *pb.Payload) {

}

type Process struct {
}

type ProgramRs struct {
	Name        string `mapstructure:"name"`
	Directory   string `mapstructure:"directory"`
	Command     string `mapstructure:"command"`
	IsAutoStart bool   `mapstructure:"isAutoStart"`
	Pid         int
	StartTime   time.Time
	StopTime    time.Time
	State       int
	StopByUser  bool
}

func SendProgramChangeRequest(programs []ProgramRs) {
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
