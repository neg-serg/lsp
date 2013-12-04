package main

import (
	"io"
	"runtime"
	"syscall"
)

const (
	blockSize = 4096
)

type file struct {
	fd      int
	name    string
	dirinfo *dirInfo // nil unless directory being read
	nepipe  int32    // number of consecutive EPIPE in Write
}

// Auxiliary information if the File describes a directory
type dirInfo struct {
	buf  []byte // buffer for directory I/O
	nbuf int    // length of buf; return value from Getdirentries
	bufp int    // location of next record in buf.
}

// syscallMode returns the syscall-specific mode bits from Go's portable mode bits.
func syscallMode(i fileMode) uint32 {
	return uint32(i)
}

func open(name string) (f *file, err error) {
	r, e := syscall.Open(name, syscall.O_RDONLY|syscall.O_CLOEXEC, syscallMode(0))
	if e != nil {
		return nil, &PathError{"open", name, e}
	}

	if syscall.O_CLOEXEC == 0 { // O_CLOEXEC not supported
		syscall.CloseOnExec(r)
	}

	return newFile(uintptr(r), name), nil
}

func newFile(fd uintptr, name string) *file {
	fdi := int(fd)
	if fdi < 0 {
		return nil
	}
	f := &file{fd: fdi, name: name}
	runtime.SetFinalizer(f, (*file).close)
	return f
}

func (file *file) close() error {
	if file == nil || file.fd < 0 {
		return syscall.EINVAL
	}
	var err error
	if e := syscall.Close(file.fd); e != nil {
		err = &PathError{"close", file.name, e}
	}
	file.fd = -1 // so it can't be closed again

	// no need for a finalizer anymore
	runtime.SetFinalizer(file, nil)
	return err
}

func (f *file) readdir(n int) (fi []*fileStat, err error) {
	dirname := f.name
	if dirname == "" {
		dirname = "."
	}
	dirname += "/"
	names, err := f.readdirnames(n)
	fi = make([]*fileStat, len(names))
	for i, filename := range names {
		fip, lerr := lstat(dirname + filename)
		if lerr != nil {
			fi[i] = &fileStat{name: filename}
			continue
		}
		fi[i] = fip
	}
	return fi, err
}
func (f *file) readdirnames(n int) (names []string, err error) {
	// If this file has no dirinfo, create one.
	if f.dirinfo == nil {
		f.dirinfo = new(dirInfo)
		// The buffer must be at least a block long.
		f.dirinfo.buf = make([]byte, blockSize)
	}
	d := f.dirinfo

	size := n
	if size <= 0 {
		size = 100
		n = -1
	}

	names = make([]string, 0, size) // Empty with room to grow.
	for n != 0 {
		// Refill the buffer if necessary
		if d.bufp >= d.nbuf {
			d.bufp = 0
			var errno error
			d.nbuf, errno = syscall.ReadDirent(f.fd, d.buf)
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
	if n >= 0 && len(names) == 0 {
		return names, io.EOF
	}
	return names, nil
}
