package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Programs []Program `mapstructure:"program"`
}

type Program struct {
	Name        string `mapstructure:"name"`
	Directory   string `mapstructure:"directory"`
	Command     string `mapstructure:"command"`
	IsAutoStart bool   `mapstructure:"isAutoStart"`
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
	return config
}
