//go:build windows

package main

import "os/exec"

func Command(program Program) *exec.Cmd {
	return exec.Command("cmd", "/C", appendPathSeparator(program.Directory)+program.Command)
}
