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

var currentPrograms []Program

func StartProgram(programs []Program) {
	currentPrograms = programs
	for _, program := range currentPrograms {
		if len(program.Directory) < 0 {
			panic(fmt.Errorf("文件目录长度必须大于0"))
		}
		if len(program.Command) < 0 {
			panic(fmt.Errorf("命令长度必须大于0"))
		}
		startProcess(&program)
	}
}

func SendProgramChangeMsg() {
	programRss := make([]program_service.ProgramRs, len(currentPrograms), len(currentPrograms))
	for index, program := range currentPrograms {
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

func Reload(programs []Program) {
	// 新增的程序
	addPrograms := getDiff(programs, currentPrograms)
	// 删除的程序
	removePrograms := getDiff(currentPrograms, programs)
	// 重新启动的程序
	restartPrograms := getRestartPrograms(programs, currentPrograms)
	for _, p := range addPrograms {
		startProcess(&p)
	}

	for _, p := range removePrograms {
		stopProcess(&p)
	}

	for _, p := range restartPrograms {
		stopProcess(&p)
		startProcess(&p)
	}

}

func getRestartPrograms(new, old []Program) []Program {
	var intersection []Program
	for _, n := range new {
		for _, o := range old {
			if n.Name == o.Name && n.Directory == o.Directory && n.Command == o.Command {
				intersection = append(intersection, n)
				break
			}
		}
	}
	return intersection
}

func getDiff(new, old []Program) []Program {
	var diff []Program
	for _, n := range new {
		found := false
		for _, o := range old {
			if n.Name == o.Name {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, n)
		}
	}
	return diff
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
