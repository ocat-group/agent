package plugin_manager

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

type State int

const (
	// Stopped the stopped state
	Stopped State = iota

	// Starting the starting state
	Starting = 10

	// Running the running state
	Running = 20

	// Backoff the backoff state
	Backoff = 30

	// Stopping the stopping state
	Stopping = 40

	// Exited the Exited state
	Exited = 100

	// Fatal the Fatal state
	Fatal = 200

	// Unknown the unknown state
	Unknown = 1000
)

type Process struct {
	process   Program
	cmd       *exec.Cmd
	startTime time.Time
	stopTime  time.Time
	state     State
	// true if process is starting
	inStart bool
	// true if the process is stopped by user
	stopByUser bool
	retryTimes *int32
	lock       sync.RWMutex
	stdin      io.WriteCloser
}

func (p Process) Start(program Program) {
	fmt.Printf("Try to start plugin_manager: %s \n", program.Name)

	command := AppendPathSeparator(program.Directory) + program.Command
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
	p.process = program
	p.startTime = time.Now()
	p.cmd = cmd
	p.state = Running
	p.inStart = true
	p.stopByUser = false
	retryTimesValue := int32(0)
	p.retryTimes = &retryTimesValue
	// 给程序赋值
	program.Process = &p
	p.checkRunning()
}

func (p *Process) checkRunning() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if !p.isRunning() {
				// 发送程序改变消息
				SendProgramChangeMsg()
				break
			}
		}
		fmt.Printf("Program exit：%s \n", p.process.Name)

		if p.process.IsAutoStart {
			// 尝试重新启动
			fmt.Printf("Try to restart plugin_manager：%s \n", p.process.Name)
			p.Start(p.process)
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

func AppendPathSeparator(path string) string {
	separator := string(os.PathSeparator)
	if !strings.HasSuffix(path, separator) {
		path += separator
	}
	return path
}
