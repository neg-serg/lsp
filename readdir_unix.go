// Implements readdir for unix without os.File

// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package main

import (
	"syscall"
)

func readdir(dirname string) ([]*fileInfo, error) {
	fd, err := open(dirname)
	if err != nil {
		return nil, err
	}
	defer syscall.Close(fd)

	dirname += "/"
	names, err := readdirnames(fd)
	fis := make([]*fileInfo, 0, len(names))
	for _, filename := range names {
		if len(filename) > 0 && filename[0] == '.' && !opts.all {
			continue
		}
		fi, _ := stat(dirname + filename)
		if fi != nil {
			fis = append(fis, fi)
		}
	}

	return fis, err
}

func open(name string) (int, error) {
	r, e := syscall.Open(name, syscall.O_RDONLY|syscall.O_CLOEXEC, 0)
	if e != nil {
		return -1, &PathError{"open", name, e}
	}

	if syscall.O_CLOEXEC == 0 { // O_CLOEXEC not supported
		syscall.CloseOnExec(r)
	}

	return r, nil
}

func readdirnames(fd int) (names []string, err error) {
	d := new(dirInfo)
	d.buf = make([]byte, blockSize)

	size := 100
	n := -1

	names = make([]string, 0, size)
	for n != 0 {
		// Refill the buffer if necessary
		if d.bufp >= d.nbuf {
			d.bufp = 0
			var errno error
			d.nbuf, errno = syscall.ReadDirent(fd, d.buf)
			if errno != nil {
				return names, NewSyscallError("readdirent", errno)
			}
			if d.nbuf <= 0 {
				break // EOF
			}
		}

		// Drain the buffer
		var nb, nc int
		nb, nc, names = syscall.ParseDirent(d.buf[d.bufp:d.nbuf], n, names)
		d.bufp += nb
		n -= nc
	}
	return names, nil
}
