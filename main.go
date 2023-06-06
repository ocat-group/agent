package main

import (
	"agent/grpc"
	"fmt"
	"github.com/jessevdk/go-flags"
	"sync"
)

// Options 命令行参数
type Options struct {
	Configuration string `short:"c" long:"configuration" description:"the configuration file"`
}

var options Options

func main() {
	// 加载命令行参数
	loadCommandLineParams()
	Reload()
}

// 加载命令行参数
func loadCommandLineParams() {
	flags.Parse(&options)
	fmt.Printf("configuration:%s \n", options.Configuration)
}

func Reload() {
	config := &Config{}
	// 加载配置文件
	config.loadConfig(&options)
	// 启动程序
	//plugin_manager.StartProgram(config.Programs)
	// todo 现阶段rpc服务只启动一次，后续要调整
	var once sync.Once
	once.Do(func() {
		grpc.StartGrpcServer(config.GrpcServerConfig)
	})
}
