package main

import (
	"agent/grpc"
	"agent/plugin_manager"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	GrpcServerConfig grpc.GrpcServerConfig    `mapstructure:"grpcServer"`
	Programs         []plugin_manager.Program `mapstructure:"program"`
}

func LoadConfig(options *Options) Config {
	viper.SetConfigFile(options.Configuration)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read config file: %s", err))
	}

	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %s", err))
	}
	printLatestConfiguration(&config)
	log.Printf("Config file loaded successfully: %s \n", options.Configuration)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s \n", e.Name)
		// 重新加载配置
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("failed to read config file: %s", err))
		}
		if err := viper.Unmarshal(&config); err != nil {
			fmt.Printf("Failed to reload config file: %s \n", err)
		} else {
			fmt.Println("Config reloaded successfully.")
			printLatestConfiguration(&config)
			Reload()
		}
	})
	// 持续运行，等待配置文件的改变
	return config
}

func printLatestConfiguration(config *Config) {
	jsonData, err := json.Marshal(config)
	if err != nil {
		fmt.Println("Failed to marshal struct to JSON:", err)
		return
	}
	fmt.Printf("Latest configuration: %v \n", string(jsonData))
}
