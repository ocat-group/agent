package main

import (
	"agent/plugin_manager"
	"fmt"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Configuration string `short:"c" long:"configuration" description:"the configuration file"`
}

var options Options

func main() {
	// 加载命令行参数
	loadCommandLineParams()
	// 加载配置文件
	config := LoadConfig(&options)
	// 启动程序
	plugin_manager.StartProgram(config.Programs)
	// 启动GRPC服务
	StartGrpcServer(config.GrpcServerConfig)
}

func loadCommandLineParams() {
	flags.NewParser(&options, flags.Default & ^flags.PrintErrors)
	fmt.Printf("configuration:%s", options.Configuration)
}
