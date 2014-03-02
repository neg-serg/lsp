package main

import (
	"syscall"
)

func typeletter(mode fileMode) string {
	switch mode & syscall.S_IFMT {
	// these are the most common, so test for them first.
	case syscall.S_IFREG:
		return cNone
	case syscall.S_IFDIR:
		return cDir
	// other letters standardized by POSIX 1003.1-2004.
	case syscall.S_IFCHR:
		return cChar
	case syscall.S_IFBLK:
		return cBlock
	case syscall.S_IFIFO:
		return cFifo
	case syscall.S_IFLNK:
		return cLink
	// other file types (though not letters) standardized by POSIX.
	case syscall.S_IFSOCK:
		return cSock
	}
	return "?"
}

// create mode strings
func strmode(buf writer, mode fileMode) {
	buf.WriteString(typeletter(mode))
	if mode&syscall.S_IRUSR != 0 {
		buf.WriteString(cRead)
	} else {
		buf.WriteString(cNone)
	}

	if mode&syscall.S_IWUSR != 0 {
		buf.WriteString(cWrite)
	} else {
		buf.WriteString(cNone)
	}

	if mode&syscall.S_ISUID != 0 {
		if mode&syscall.S_IXUSR != 0 {
			buf.WriteString(cUIDExec)
		} else {
			buf.WriteString(cUID)
		}
	} else if mode&syscall.S_IXUSR != 0 {
		buf.WriteString(cExec)
	} else {
		buf.WriteString(cNone)
	}

	if mode&syscall.S_IRGRP != 0 {
		buf.WriteString(cRead)
	} else {
		buf.WriteString(cNone)
	}

	if mode&syscall.S_IWGRP != 0 {
		buf.WriteString(cWrite)
	} else {
		buf.WriteString(cNone)
	}

	if mode&syscall.S_ISGID != 0 {
		if mode&syscall.S_IXGRP != 0 {
			buf.WriteString(cUIDExec)
		} else {
			buf.WriteString(cUID)
		}
	} else if mode&syscall.S_IXGRP != 0 {
		buf.WriteString(cExec)
	} else {
		buf.WriteString(cNone)
	}

	if mode&syscall.S_IROTH != 0 {
		buf.WriteString(cRead)
	} else {
		buf.WriteString(cNone)
	}

	if mode&syscall.S_IWOTH != 0 {
		buf.WriteString(cWrite)
	} else {
		buf.WriteString(cNone)
	}

	if mode&syscall.S_ISVTX != 0 {
		if mode&syscall.S_IXOTH != 0 {
			buf.WriteString(cSticky)
		} else {
			buf.WriteString(cStickyO)
		}
	} else if mode&syscall.S_IXOTH != 0 {
		buf.WriteString(cExec)
	} else {
		buf.WriteString(cNone)
	}

	buf.WriteString("\033[0m")
}
