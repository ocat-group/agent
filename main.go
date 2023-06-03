package main

func main() {
	// 加载配置文件
	config := LoadConfig()
	// 启动程序
	StartProgram(config.Programs)
	// 启动GRPC服务
	StartGrpcServer(config.GrpcServerConfig)
}
