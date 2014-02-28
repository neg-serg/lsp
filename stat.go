package main

import (
	"syscall"
	"time"
)

type fileMode uint64

type fileInfo struct {
	name string
	size int64
	mode fileMode
	time int64
}

func (fi *fileInfo) isDir() bool {
	return fi.mode&syscall.S_IFMT == syscall.S_IFDIR
}

func (m fileMode) isDir() bool {
	return m&syscall.S_IFMT == syscall.S_IFDIR
}

func (m fileMode) isRegular() bool {
	return m&syscall.S_IFMT == syscall.S_IFREG
}

// get info about file/directory name
func ls(name string) (fileList, bool, error) {
	fi, err := stat(name)
	if err != nil {
		return nil, false, err
	}
	if fi.isDir() {
		fis, err := readdir(name)
		if *all {
			return fis, true, err
		}
		filtered := make([]*fileInfo, 0, len(fis))
		for _, fi := range fis {
			if len(fi.name) > 0 && fi.name[0] == '.' {
				continue
			}
			filtered = append(filtered, fi)
		}
		return filtered, true, err
	}
	return []*fileInfo{fi}, false, nil
}

// stat returns a fileInfo describing the named file
func stat(name string) (fi *fileInfo, err error) {
	var stat syscall.Stat_t
	err = syscall.Stat(name, &stat)
	if err != nil {
		return nil, &PathError{"stat", name, err}
	}
	return fileInfoFromStat(&stat, name), nil
}

// lstat returns a fileInfo describing the named file. If the file is a
// symbolic link, the returned *fileInfo describes the symbolic link.
// lstat makes no attempt to follow the link.
func lstat(name string) (fi *fileInfo, err error) {
	var stat syscall.Stat_t
	err = syscall.Lstat(name, &stat)
	if err != nil {
		return nil, &PathError{"lstat", name, err}
	}
	return fileInfoFromStat(&stat, name), nil
}

func fileInfoFromStat(st *syscall.Stat_t, name string) *fileInfo {
	f := &fileInfo{
		name: basename(name),
		size: int64(st.Size),
		mode: fileMode(st.Mode),
	}
	var t syscall.Timespec
	if *ctime {
		t = st.Ctim
	} else {
		t = st.Mtim
	}
	f.time = int64(t.Sec)*1e9 + int64(t.Nsec)
	return f
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
			return string(b[0:n]), nil
		}
	}
}
