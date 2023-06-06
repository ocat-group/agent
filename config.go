package main

import (
	"agent/grpc"
	"agent/plugin_manager"
	"agent/util"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"sync"
)

var config Config

type Config struct {
	GrpcServerConfig grpc.ServerConfig        `mapstructure:"grpcServer"`
	Programs         []plugin_manager.Program `mapstructure:"program"`
	Lock             *sync.Mutex
	Md5              string
}

var configInit sync.Once
var configListen sync.Once

func (c *Config) loadConfig(options *Options) Config {
	viper.SetConfigFile(options.Configuration)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read config file: %v", err))
	}
	// 初始化config
	configInit.Do(func() {
		config = Config{Lock: &sync.Mutex{}}
	})
	// 加载config
	config.Lock.Lock()
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %s", err))
	}
	config.Lock.Unlock()
	// 监听配置
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		currentMd5, err := util.ReadFileMd5(options.Configuration)
		if err != nil {
			fmt.Println("Failed to read configuration file.")
		}
		if config.Md5 == currentMd5 {
			fmt.Println("The configuration file has not changed, skip this processing.")
			return
		}
		config.Md5 = currentMd5
		fmt.Printf("Config file changed: %s \n", e.Name)

		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("failed to read config file: %s", err))
		}
		if err := viper.Unmarshal(&config); err != nil {
			fmt.Printf("Failed to reload config file: %s \n", err)
		} else {
			fmt.Println("Config reloaded successfully.")
			config.printLatestConfig()
			Reload()
		}
	})
	return config
}

func (c *Config) printLatestConfig() {
	jsonData, err := json.Marshal(c)
	if err != nil {
		fmt.Println("Failed to marshal struct to JSON.", err)
		return
	}
	fmt.Printf("Latest configuration: %v \n", string(jsonData))
}
