package main

import (
	"agent/grpc"
	"agent/plugin_manager"
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
	// 重载插件
	plugin_manager.Reload(config.Programs)
	// todo grpc目前先不考虑端口的改变，因为端口改变会导致所有长链接断开，
	// todo 服务端需要获取最新的端口进行重连，这里是否有意义，待考虑，所以暂时不进行实现
	var once sync.Once
	once.Do(func() {
		grpc.StartGrpcServer(config.GrpcServerConfig)
	})
}
