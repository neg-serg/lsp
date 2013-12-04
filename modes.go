package main

import (
	"bytes"
	"syscall"
)

const (
	cESC = "\033["
	cEnd = cESC + "0m"

	cChar     = cESC + "0m" + "c"
	cDev      = cESC + "0m" + "b"
	cDir      = cESC + "38;5;2;1m" + "d" + cEnd
	cExec     = cESC + "38;5;131m" + "x"
	cFifo     = cESC + "0m" + "p"
	cLink     = cESC + "38;5;220;1m" + "l" + cEnd
	cNone     = cESC + "38;5;0m" + "—"
	cRead     = cESC + "38;5;2m" + "r"
	cRes      = cESC + "38;5;220m" + "t"
	cResOther = cESC + "38;5;220;1m" + "T" + cEnd
	cUID      = cESC + "38;5;220m" + "S"
	cUIDExec  = cESC + "38;5;161m" + "s"
	cSock     = cESC + "38;5;161m" + "s"
	cWrite    = cESC + "38;5;216m" + "w"
)

func typeletter(mode fileMode) string {
	switch {
	// These are the most common, so test for them first.
	case mode.isRegular():
		return cNone
	case mode.isDir():
		return cDir
	}

	switch mode & syscall.S_IFMT {
	// Other letters standardized by POSIX 1003.1-2004.
	case syscall.S_IFCHR:
		return cChar
	case syscall.S_IFBLK:
		return cDev
	case syscall.S_IFIFO:
		return cFifo
	case syscall.S_IFLNK:
		return cLink

	// Other file types (though not letters) standardized by POSIX.
	case syscall.S_IFSOCK:
		return cSock
	}
	return "?"
}

var buf bytes.Buffer

/* Like filemodestring, but rely only on MODE.  */
func strmode(mode fileMode) string {
	buf.Reset()
	buf.WriteString(typeletter(mode))

	if mode&modeIRUSR != 0 {
		buf.WriteString(cRead)
	} else {
		buf.WriteString(cNone)
	}

	if mode&modeIWUSR != 0 {
		buf.WriteString(cWrite)
	} else {
		buf.WriteString(cNone)
	}

	if mode&syscall.S_ISUID != 0 {
		if mode&modeIXUSR != 0 {
			buf.WriteString(cUIDExec)
		} else {
			buf.WriteString(cUID)
		}
	} else if mode&modeIXUSR != 0 {
		buf.WriteString(cExec)
	} else {
		buf.WriteString(cNone)
	}

	if mode&modeIRGRP != 0 {
		buf.WriteString(cRead)
	} else {
		buf.WriteString(cNone)
	}

	if mode&modeIWGRP != 0 {
		buf.WriteString(cWrite)
	} else {
		buf.WriteString(cNone)
	}

	if mode&syscall.S_ISGID != 0 {
		if mode&modeIXGRP != 0 {
			buf.WriteString(cUIDExec)
		} else {
			buf.WriteString(cUID)
		}
	} else if mode&modeIXGRP != 0 {
		buf.WriteString(cExec)
	} else {
		buf.WriteString(cNone)
	}

	if mode&modeIROTH != 0 {
		buf.WriteString(cRead)
	} else {
		buf.WriteString(cNone)
	}

	if mode&modeIWOTH != 0 {
		buf.WriteString(cWrite)
	} else {
		buf.WriteString(cNone)
	}

	if mode&syscall.S_ISVTX != 0 {
		if mode&modeIXOTH != 0 {
			buf.WriteString(cRes)
		} else {
			buf.WriteString(cResOther)
		}
	} else if mode&modeIXOTH != 0 {
		buf.WriteString(cExec)
	} else {
		buf.WriteString(cNone)
	}

	buf.WriteString("\033[0m")
	return buf.String()
}

/*
func strmode(m os.FileMode) string {
	const str = "dalTLDpSugct"
	var buf [32]byte // mode is uint32.
	w := 0
	for i, c := range str {
		if m&(1<<uint(32-1-i)) != 0 {
			buf[w] = byte(c)
			w++
		}
	}
	if w == 0 {
		buf[w] = '-'
		w++
	}
	const rwx = "rwxrwxrwx"
	for i, c := range rwx {
		if m&(1<<uint(9-1-i)) != 0 {
			buf[w] = byte(c)
		} else {
			buf[w] = '-'
		}
		w++
	}
	return string(buf[:w])
}
*/