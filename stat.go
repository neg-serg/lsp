package main

import (
	"strings"
	"syscall"
	"time"
)

type fileMode uint64

type fileInfo struct {
	name     string
	size     int64
	mode     fileMode
	time     int64
	linkok   bool
	linkname string
	linkmode fileMode
}

// get info about file/directory name
func ls(name string) (fileList, error) {
	fi, err := stat(name)
	if err != nil {
		return nil, err
	}
	if fi.mode&syscall.S_IFMT == syscall.S_IFDIR {
		return readdir(name)
	}
	return []fileInfo{*fi}, nil
}

func gettime(st *syscall.Stat_t) int64 {
	var t syscall.Timespec
	if args.ctime {
		t = st.Ctim
	} else {
		t = st.Mtim
	}
	return int64(t.Sec)*1e9 + int64(t.Nsec)
}

// stat returns a fileInfo describing the named file
func stat(name string) (*fileInfo, error) {
	var stat syscall.Stat_t
	err := syscall.Lstat(name, &stat)
	if err != nil {
		return nil, &PathError{"stat", name, err}
	}

	fi := &fileInfo{
		name:   basename(name),
		size:   int64(stat.Size),
		mode:   fileMode(stat.Mode),
		time:   gettime(&stat),
		linkok: true,
	}

	if fi.mode&syscall.S_IFMT == syscall.S_IFLNK {
		fi.linkname, err = readlink(name)
		if err != nil {
			fi.linkok = false
			return fi, nil
		}
		err = syscall.Stat(name, &stat)
		if err != nil {
			fi.linkok = false
			return fi, nil
		}
		fi.linkmode = fileMode(stat.Mode)
	}

	return fi, nil
}

// basename the leading directory name from path name
func basename(name string) string {
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '/' {
			return name[i+1:]
		}
	}
	return name
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

func readlink(name string) (string, error) {
	for len := 128; ; len *= 2 {
		b := make([]byte, len)
		n, e := syscall.Readlink(name, b)
		if e != nil {
			return "", &PathError{"readlink", name, e}
		}
		if n < len {
			return strings.TrimRight(string(b[0:n]), "/"), nil
		}
	}
}
