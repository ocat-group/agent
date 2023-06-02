package main

var config Config

func main() {
	config := LoadConfig()
	for _, program := range config.Programs {
		Start(program)
	}
}
