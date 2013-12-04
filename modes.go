package main

import (
	"bytes"
	"syscall"
)

const (
	cESC = "\033["
	cEnd = cESC + "0m"

	cNone = cESC + "38;5;0m" + "—"

	cChar     = cESC + "0m" + "c"
	cDev      = cESC + "0m" + "b"
	cDir      = cESC + "38;5;2;1m" + "d" + cEnd
	cExec     = cESC + "38;5;131m" + "x"
	cFifo     = cESC + "0m" + "p"
	cLink     = cESC + "38;5;220;1m" + "l" + cEnd
	cRead     = cESC + "38;5;2m" + "r"
	cRes      = cESC + "38;5;220m" + "t"
	cResOther = cESC + "38;5;220;1m" + "T" + cEnd
	cSock     = cESC + "38;5;161m" + "s"
	cUID      = cESC + "38;5;220m" + "S"
	cUIDExec  = cESC + "38;5;161m" + "s"
	cWrite    = cESC + "38;5;216m" + "w"
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
		return cDev
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

var buf bytes.Buffer

// create mode strings
func strmode(mode fileMode) string {
	buf.Reset()
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
			buf.WriteString(cRes)
		} else {
			buf.WriteString(cResOther)
		}
	} else if mode&syscall.S_IXOTH != 0 {
		buf.WriteString(cExec)
	} else {
		buf.WriteString(cNone)
	}

	buf.WriteString("\033[0m")
	return buf.String()
}
