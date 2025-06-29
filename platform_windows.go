// platform_windows.go
//go:build windows

package main

import (
	"golang.org/x/sys/windows"
	"syscall"
)

func launchGamePlatform(exePath string, workDir string) error {
	verbPtr, _ := syscall.UTF16PtrFromString("runas")
	exePtr, _ := syscall.UTF16PtrFromString(exePath)
	argPtr, _ := syscall.UTF16PtrFromString("")
	cwdPtr, _ := syscall.UTF16PtrFromString(workDir)
	return windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, 1)
}
