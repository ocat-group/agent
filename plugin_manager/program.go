package plugin_manager

import (
	program_service "agent/service"
	"agent/util"
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"
)

const (
	// Stopped the stopped state
	Stopped = 0

	// Starting the starting state
	Starting = 10

	// Running the running state
	Running = 20

	// Backoff the backoff state
	Backoff = 30

	// Stopping the stopping state
	Stopping = 40

	// Exited the Exited state
	Exited = 50

	// Fatal the Fatal state
	Fatal = 60

	// Unknown the unknown state
	Unknown = 70
)

type Program struct {
	Name            string `mapstructure:"name"`
	Directory       string `mapstructure:"directory"`
	Command         string `mapstructure:"command"`
	IsAutoStart     bool   `mapstructure:"isAutoStart"`
	MaxRestartCount int    `mapstructure:"MaxRestartCount"`
	Process         *Process
}

type Process struct {
	program   *Program
	cmd       *exec.Cmd
	startTime time.Time
	stopTime  time.Time
	state     int
	// true if process is starting
	inStart bool
	// true if the process is stopped by user
	stopByUser bool
	retryTimes int
}

var currentPrograms []*Program

func Reload(programs []*Program) {
	checkAndRemove(programs)
	addPrograms := computesAddPrograms(programs, currentPrograms)
	removePrograms := computesRemovePrograms(currentPrograms, programs)
	restartPrograms := computesRestartPrograms(programs, currentPrograms)
	for _, p := range addPrograms {
		p.start()
	}
	for _, p := range removePrograms {
		p.stop()
	}
	for _, p := range restartPrograms {
		p.stop()
		p.start()
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

func (p *Program) start() {
	fmt.Printf("Try to start Plugin: %s \n", p.Name)
	cmd := p.startProcess()
	p.updateProgramToStart(cmd)
	p.listenRunningStatus()
}

func (p *Program) listenRunningStatus() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if !p.isRunning() {
				p.updateProgramToStop()
				// 发送程序改变消息
				SendProgramChangeMsg()
				break
			}
		}
		fmt.Printf("Program exit：%s \n", p.Name)

		if p.IsAutoStart {
			// 尝试重新启动
			fmt.Printf("Try to restart Plugin：%s \n", p.Name)
			p.start()
			SendProgramChangeMsg()
		}
	}()
}

func (p *Program) isRunning() bool {
	if p.Process.cmd != nil && p.Process.cmd.Process != nil {
		if runtime.GOOS == "windows" {
			exists, err := process.PidExists(int32(p.Process.cmd.Process.Pid))
			return exists && err == nil
		}
		return p.Process.cmd.Process.Signal(syscall.Signal(0)) == nil
	}
	return false
}

func (p *Program) updateProgramToStart(cmd *exec.Cmd) {
	if p.Process != nil {
		p.Process.startTime = time.Now()
		p.Process.state = Running
		p.Process.inStart = true
		p.Process.retryTimes++
		return
	}
	// 第一次启动，创建Process
	process := new(Process)
	process.cmd = cmd
	process.startTime = time.Now()
	process.stopTime = time.Time{}
	process.state = Running
	process.inStart = true
	process.stopByUser = false
	process.retryTimes = 0
	// 给进程赋值
	p.Process = process
}

func (p *Program) updateProgramToStop() {
	p.Process.cmd = nil
	p.Process.startTime = time.Time{}
	p.Process.stopTime = time.Now()
	p.Process.state = Stopped
	p.Process.inStart = false
	p.Process.stopByUser = false
}

func (p *Program) check() error {
	if len(p.Directory) < 0 {
		return errors.New("the file directory length must be greater than 0")
	}
	if len(p.Command) < 0 {
		return errors.New("command length must be greater than 0")
	}
	return nil
}

func checkAndRemove(programs []*Program) {
	for i := 0; i < len(programs); i++ {
		if err := programs[i].check(); err != nil {
			programs = append(programs[:i], programs[i+1:]...)
			i--
		}
	}
}

func computesRestartPrograms(new, old []*Program) []*Program {
	var restart []*Program
	for _, n := range new {
		found := false
		for _, o := range old {
			if n.Name == o.Name && (n.Directory != o.Directory || n.Command != o.Command) {
				found = true
				break
			}
		}
		if found {
			restart = append(restart, n)
		}
	}
	return restart
}

func computesAddPrograms(new, old []*Program) []*Program {
	return computesDifference(new, old)
}

func computesRemovePrograms(new, old []*Program) []*Program {
	return computesDifference(old, new)
}

func computesDifference(new, old []*Program) []*Program {
	if new == nil {
		return nil
	} else if old == nil {
		return new
	}
	var diff []*Program
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

func (p *Program) stop() {
	err := p.Process.cmd.Process.Kill()
	if err != nil {
		log.Printf("停止进程失败: %s \n", err)
	}
}

func (p *Program) startProcess() *exec.Cmd {
	completeCommand := util.AppendPathSeparator(p.Directory) + p.Command
	cmd := exec.Command(completeCommand)
	cmd.Dir = p.Directory
	// 设置标准输出和标准错误输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// 启动命令
	err := cmd.Start()
	if err != nil {
		panic(fmt.Errorf("启动命令时出错: %s", err))
	}
	return cmd
}
