package main

import (
	"agent/grpc"
	"agent/plugin_manager"
	"fmt"
	"github.com/jessevdk/go-flags"
)

// Options 命令行参数
type Options struct {
	Configuration string `short:"c" long:"configuration" description:"the configuration file"`
}

var options Options

var config Config

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
	// 加载配置文件
	config := LoadConfig(&options)
	// 启动程序
	plugin_manager.StartProgram(config.Programs)
	// 启动GRPC服务
	grpc.StartGrpcServer(config.GrpcServerConfig)
}
