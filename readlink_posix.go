// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris windows

package main

import "syscall"

func readlink(name string) (string, error) {
	for len := 128; ; len *= 2 {
		b := make([]byte, len)
		n, e := syscall.Readlink(name, b)
		if e != nil {
			return "", &PathError{"readlink", name, e}
		}
		if n < len {
			return string(cleanRight(b[0:n])), nil
		}
	}
}
