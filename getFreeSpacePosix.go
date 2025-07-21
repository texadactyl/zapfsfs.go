//go:build !windows
// +build !windows

package main

import (
	"golang.org/x/sys/unix"
)

func getFreeSpace(path string) (uint64, error) {
	var stat unix.Statfs_t
	err := unix.Statfs(path, &stat)
	if err != nil {
		return 0, err
	}
	return stat.Bavail * uint64(stat.Bsize), nil
}
