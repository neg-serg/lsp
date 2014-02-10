package main

import (
	"syscall"
)

const (
	blockSize = 4096
)

// Auxiliary information if the file describes a directory
type dirInfo struct {
	buf  []byte // buffer for directory I/O
	nbuf int    // length of buf; return value from Getdirentries
	bufp int    // location of next record in buf.
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

func readdir(fname string) (fi []*fileInfo, err error) {
	fd, err := open(fname)
	if err != nil {
		return nil, err
	}
	defer syscall.Close(fd)

	dirname := fname
	if dirname == "" {
		dirname = "."
	}
	dirname += "/"
	names, err := readdirnames(fd)
	fi = make([]*fileInfo, len(names))
	for i, filename := range names {
		fip, lerr := lstat(dirname + filename)
		if lerr != nil {
			fi[i] = &fileInfo{name: filename}
			continue
		}
		fi[i] = fip
	}
	return fi, err
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
