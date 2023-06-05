package main

import (
	"agent/plugin_manager"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	GrpcServerConfig GrpcServerConfig         `mapstructure:"grpcServer"`
	Programs         []plugin_manager.Program `mapstructure:"plugin_manager"`
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
	log.Printf("Config file loaded successfully: %s", options.Configuration)
	return config
}
