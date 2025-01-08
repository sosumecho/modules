//go:build darwin || freebsd || linux
// +build darwin freebsd linux

package utils

import "syscall"

const (
	SIGUSR1 = syscall.SIGUSR1
	SIGUSR2 = syscall.SIGUSR2

	LOCK_EX = syscall.LOCK_EX
	LOCK_NB = syscall.LOCK_NB
)

// Flock 文件加锁
func Flock(fd int, how int) error {
	return syscall.Flock(fd, how)
}
