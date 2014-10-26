package main

import (
	"os"
	"syscall"
	"unsafe"
)

func isTTY(f *os.File) bool {
	var t syscall.Termios
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), ioctlGetTermios, uintptr(unsafe.Pointer(&t)))
	return err == 0
}
