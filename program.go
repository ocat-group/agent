package main

import (
	"fmt"
)

type Program struct {
	Name        string `mapstructure:"name"`
	Directory   string `mapstructure:"directory"`
	Command     string `mapstructure:"command"`
	IsAutoStart bool   `mapstructure:"isAutoStart"`
}

func StartProgram(programs []Program) {
	for _, program := range programs {
		if len(program.Directory) < 0 {
			panic(fmt.Errorf("文件目录长度必须大于0"))
		}
		if len(program.Command) < 0 {
			panic(fmt.Errorf("命令长度必须大于0"))
		}
		p := new(Process)
		p.Start(program)
	}
}
