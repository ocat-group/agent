package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestStartProgram(t *testing.T) {
	cmd := exec.Command("cmd", "/C", "D:\\DevSoftware\\Scripte\\test.bat qwe cas")
	cmd.Dir = "D:\\DevSoftware\\Scripte"

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

func TestAppendPathSeparator(t *testing.T) {
	path := appendPathSeparator("D:\\DevSoftware\\Scripte")
	fmt.Println(path)
}
