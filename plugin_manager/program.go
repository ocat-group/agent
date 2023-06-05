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

func SendProgramChangeMsg() {
	programRss := make([]program_service.ProgramRs, len(programs), len(programs))
	for index, program := range programs {
		programRss[index] = program_service.ProgramRs{Name: program.Name,
			Directory:   program.Directory,
			Command:     program.Command,
			IsAutoStart: program.IsAutoStart,
			Pid:         program.Process.cmd.Process.Pid,
			StartTime:   program.Process.startTime,
			StopTime:    program.Process.stopTime,
			State:       program.Process.state,
			StopByUser:  program.Process.stopByUser,
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
