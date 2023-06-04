package plugin_manager

import (
	program_service "agent/service"
	"fmt"
	"time"
)

type Program struct {
	Name            string `mapstructure:"name"`
	Directory       string `mapstructure:"directory"`
	Command         string `mapstructure:"command"`
	IsAutoStart     bool   `mapstructure:"isAutoStart"`
	MaxRestartCount int    `mapstructure:"MaxRestartCount"`
	Process         *Process
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
		startProcess(program)
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

func checkRunning(program Program) {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if !program.Process.isRunning() {
				processStop(program.Process)
				// 发送程序改变消息
				SendProgramChangeMsg()
				break
			}
		}
		fmt.Printf("Program exit：%s \n", program.Process.process.Name)

		if program.Process.process.IsAutoStart {
			if program.MaxRestartCount <= program.Process.retryTimes {
				fmt.Println("已经达到最大重启次数，不再进行重启")
				return
			}
			// 尝试重新启动
			fmt.Printf("Try to restart Plugin：%s \n", program.Process.process.Name)
			startProcess(program.Process.process)
			SendProgramChangeMsg()
		}
	}()
}
