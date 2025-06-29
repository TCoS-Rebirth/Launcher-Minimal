// platform_unix.go
//go:build !windows

package main

import (
	"os/exec"
)

func launchGamePlatform(exePath string, workDir string) error {
	cmd := exec.Command(exePath)
	cmd.Dir = workDir
	return cmd.Start()
}
