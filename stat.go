package main

import (
	"syscall"
	"time"
)

// get info about file/directory name
func ls(name string) ([]*fileStat, error) {
	fi, err := stat(name)
	if err != nil {
		return nil, err
	}
	if fi.isDir() {
		f, err := open(name)
		if err != nil {
			return nil, err
		}
		defer f.close()
		fis, err := f.readdir(0)
		if *all {
			return fis, err
		}
		filtered := make([]*fileStat, 0, len(fis))
		for _, fi := range fis {
			if len(fi.name) > 0 && fi.name[0] == '.' {
				continue
			}
			filtered = append(filtered, fi)
		}
		return filtered, err
	}
	return []*fileStat{fi}, nil
}

type fileMode uint64

type fileStat struct {
	name    string
	size    int64
	mode    fileMode
	cTime   time.Time
	modTime time.Time
}

// Stat returns a fileStat describing the named file.
// If there is an error, it will be of type *PathError.
func stat(name string) (fi *fileStat, err error) {
	var stat syscall.Stat_t
	err = syscall.Stat(name, &stat)
	if err != nil {
		return nil, &PathError{"stat", name, err}
	}
	return fileStatFromStat(&stat, name), nil
}

// Lstat returns a FileInfo describing the named file.
// If the file is a symbolic link, the returned FileInfo
// describes the symbolic link.  Lstat makes no attempt to follow the link.
// If there is an error, it will be of type *PathError.
func lstat(name string) (fi *fileStat, err error) {
	var stat syscall.Stat_t
	err = syscall.Lstat(name, &stat)
	if err != nil {
		return nil, &PathError{"lstat", name, err}
	}
	return fileStatFromStat(&stat, name), nil
}

func fileStatFromStat(st *syscall.Stat_t, name string) *fileStat {
	fs := &fileStat{
		name:    basename(name),
		size:    int64(st.Size),
		modTime: timespecToTime(st.Mtim),
		cTime:   timespecToTime(st.Ctim),
	}
	fs.mode = fileMode(st.Mode)
	return fs
}

// basename removes trailing slashes and the leading directory name from path name
func basename(name string) string {
	i := len(name) - 1
	// Remove trailing slashes
	for ; i > 0 && name[i] == '/'; i-- {
		name = name[:i]
	}
	// Remove leading directory name
	for i--; i >= 0; i-- {
		if name[i] == '/' {
			name = name[i+1:]
			break
		}
	}

	return name
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

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

func (fi *fileStat) isDir() bool {
	return fi.mode&syscall.S_IFMT == syscall.S_IFDIR
}

func (m fileMode) isDir() bool {
	return m&syscall.S_IFMT == syscall.S_IFDIR
}

func (m fileMode) isRegular() bool {
	return m&syscall.S_IFMT == syscall.S_IFREG
}

func (m fileMode) Perm() fileMode {
	return m
}
