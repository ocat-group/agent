package main

import (
	"fmt"
	"os"
	"strings"
)

func Start(programs []Program) {
	for _, program := range programs {
		if len(program.Directory) < 0 {
			panic(fmt.Errorf("文件目录长度必须大于0"))
		}
		if len(program.Command) < 0 {
			panic(fmt.Errorf("命令长度必须大于0"))
		}
		p := new(Process)
		p.Start(program)
	}
}

func appendPathSeparator(path string) string {
	separator := string(os.PathSeparator)
	if !strings.HasSuffix(path, separator) {
		path += separator
	}
	return path
}
