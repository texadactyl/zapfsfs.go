//go:build windows
// +build windows

package main

import (
	"golang.org/x/sys/windows"
	"syscall"
)

func getFreeSpace(path string) (uint64, error) {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}

	var avail, total, free uint64
	err = windows.GetDiskFreeSpaceEx(pathPtr, &avail, &total, &free)
	if err != nil {
		return 0, err
	}
	return avail, nil
}
