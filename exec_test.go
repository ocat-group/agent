package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestCommand(t *testing.T) {
	cmd := exec.Command("E:\\GolangProjects\\agent\\script\\test_script.bat")
	cmd.Dir = "E:\\GolangProjects\\agent\\script\\"

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()

	if err != nil {
		panic(fmt.Errorf("启动命令时出错: %s", err))
	}
	err = cmd.Wait()
	if err != nil {
		panic(fmt.Errorf("等待命令完成时出错：%s", err))
	}
}
