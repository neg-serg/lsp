package main

import "syscall"

type fileMode uint64

type fileInfo struct {
	name     sufIndexed
	size     int64
	mode     fileMode
	time     int64
	linkok   bool
	linkname sufIndexed
	linkmode fileMode
}

// get info about file/directory name
func ls(name string) (fileList, error) {
	fi := &fileInfo{}
	err := stat(name, fi)
	if err != nil {
		return nil, err
	}
	if fi.mode&syscall.S_IFMT == syscall.S_IFDIR {
		return readdir(name)
	}
	return []*fileInfo{fi}, nil
}

// stat returns a fileInfo describing the named file
func stat(name string, out *fileInfo) error {
	var stat syscall.Stat_t
	err := syscall.Lstat(name, &stat)
	if err != nil {
		return &PathError{"stat", name, err}
	}

	*out = fileInfo{
		name:   newSufIndexed(basename(name)),
		size:   int64(stat.Size),
		mode:   fileMode(stat.Mode),
		time:   gettime(&stat),
		linkok: true,
	}

	if out.mode&syscall.S_IFMT == syscall.S_IFLNK {
		ln, err := readlink(name)
		if err != nil {
			out.linkok = false
			return nil
		}
		out.linkname = newSufIndexed(ln)
		err = syscall.Stat(name, &stat)
		if err != nil {
			out.linkok = false
			return nil
		}
		out.linkmode = fileMode(stat.Mode)
	}

	return nil
}
