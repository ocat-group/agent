package main

import (
	"fmt"
	"os"
	"strings"
)

func Start(program Program) {
	cmd := Command(program)
	cmd.Dir = program.Directory
	// 设置标准输出和标准错误输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// 启动命令
	err := cmd.Start()
	if err != nil {
		panic(fmt.Errorf("启动命令时出错: %s", err))
	}
}

func appendPathSeparator(path string) string {
	separator := string(os.PathSeparator)
	if !strings.HasSuffix(path, separator) {
		path += separator
	}
	return path
}
