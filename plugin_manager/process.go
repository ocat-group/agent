package plugin_manager

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
)

type State int

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

type Process struct {
	process   Program
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

func startProcess(program Program) {
	fmt.Printf("Try to start Plugin: %s \n", program.Name)

	command := appendPathSeparator(program.Directory) + program.Command
	cmd := exec.Command(command)
	cmd.Dir = program.Directory
	// 设置标准输出和标准错误输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// 启动命令
	err := cmd.Start()
	if err != nil {
		panic(fmt.Errorf("启动命令时出错: %s", err))
	}
	// 创建Process
	p := createProcess(program, cmd)
	p.listenRunningStatus()
}

func createProcess(program Program, cmd *exec.Cmd) *Process {
	// 代表第一次启动
	if program.Process != nil {
		program.Process.startTime = time.Now()
		program.Process.state = Running
		program.Process.inStart = true
		program.Process.retryTimes++
		return program.Process
	}
	// 创建Process
	p := new(Process)
	p.process = program
	p.cmd = cmd
	p.startTime = time.Now()
	//p.stopTime
	p.state = Running
	p.inStart = true
	p.stopByUser = false
	p.retryTimes = 0
	// 给进程赋值
	program.Process = p
	return p
}

func processStop(p *Process) {
	p.stopTime = time.Now()
	p.state = Stopped
	p.inStart = false
}

func (p *Process) listenRunningStatus() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if !p.isRunning() {
				processStop(p)
				// 发送程序改变消息
				SendProgramChangeMsg()
				break
			}
		}
		fmt.Printf("Program exit：%s \n", p.process.Name)

		if p.process.IsAutoStart {
			// 尝试重新启动
			fmt.Printf("Try to restart Plugin：%s \n", p.process.Name)
			startProcess(p.process)
			SendProgramChangeMsg()
		}
	}()
}

// check if the process is running or not
//
func (p *Process) isRunning() bool {
	if p.cmd != nil && p.cmd.Process != nil {
		if runtime.GOOS == "windows" {
			exists, err := process.PidExists(int32(p.cmd.Process.Pid))
			return exists && err == nil
		}
		return p.cmd.Process.Signal(syscall.Signal(0)) == nil
	}
	return false
}

func appendPathSeparator(path string) string {
	separator := string(os.PathSeparator)
	if !strings.HasSuffix(path, separator) {
		path += separator
	}
	return path
}
