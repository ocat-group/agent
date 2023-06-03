package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	GrpcServerConfig GrpcServerConfig `mapstructure:"grpcServer"`
	Programs         []Program        `mapstructure:"program"`
}

func LoadConfig() Config {
	viper.SetConfigFile("config.yml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read config file: %s", err))
	}

	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %s", err))
	}
	log.Printf("Config file loaded successfully: %s", "config.yml")
	return config
}
