package plugin_manager

import (
	program_service "agent/service"
	"fmt"
	"time"
)

type Program struct {
	Name        string `mapstructure:"name"`
	Directory   string `mapstructure:"directory"`
	Command     string `mapstructure:"command"`
	IsAutoStart bool   `mapstructure:"isAutoStart"`
	Process     *Process
}

var programs []Program

func StartProgram(programs []Program) {
	for _, program := range programs {
		if len(program.Directory) < 0 {
			panic(fmt.Errorf("文件目录长度必须大于0"))
		}
		if len(program.Command) < 0 {
			panic(fmt.Errorf("命令长度必须大于0"))
		}
		Process{}.Start(program)
	}
}

type ProgramRs struct {
	name        string `mapstructure:"name"`
	directory   string `mapstructure:"directory"`
	command     string `mapstructure:"command"`
	isAutoStart bool   `mapstructure:"isAutoStart"`
	pid         int
	startTime   time.Time
	stopTime    time.Time
	state       State
	stopByUser  bool
}

func SendProgramChangeMsg() {
	programRss := make([]ProgramRs, len(programs), len(programs))
	for index, program := range programs {
		programRss[index] = ProgramRs{name: program.Name,
			directory:   program.Directory,
			command:     program.Command,
			isAutoStart: program.IsAutoStart,
			pid:         program.Process.cmd.Process.Pid,
			startTime:   program.Process.startTime,
			stopTime:    program.Process.stopTime,
			state:       program.Process.state,
			stopByUser:  program.Process.stopByUser,
		}
	}
	program_service.SendProgramChangeRequest(programRss)
}
