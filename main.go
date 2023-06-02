package main

import "time"

var config Config

func main() {
	config := LoadConfig()
	Start(config.Programs)
	time.Sleep(1000 * time.Second)
}
