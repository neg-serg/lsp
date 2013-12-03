package main

import (
	"bytes"
	"os"
)

const (
	ESC  = "\033["
	cEnd = ESC + "0m"

	cChar     = ESC + "0m" + "c"
	cDev      = ESC + "0m" + "b"
	cDir      = ESC + "38;5;2;1m" + "d" + cEnd
	cExec     = ESC + "38;5;131m" + "x"
	cFifo     = ESC + "0m" + "p"
	cLink     = ESC + "38;5;220;1m" + "l" + cEnd
	cNone     = ESC + "38;5;0m" + "—"
	cRead     = ESC + "38;5;2m" + "r"
	cRes      = ESC + "38;5;220m" + "t"
	cResOther = ESC + "38;5;220;1m" + "T" + cEnd
	cUid      = ESC + "38;5;220m" + "S"
	cUidExec  = ESC + "38;5;161m" + "s"
	cSock     = ESC + "38;5;161m" + "s"
	cWrite    = ESC + "38;5;216m" + "w"
)

func typeletter(mode os.FileMode) string {
	switch {
	// These are the most common, so test for them first.
	case mode.IsRegular():
		return cNone
	case mode.IsDir():
		return cDir

	// Other letters standardized by POSIX 1003.1-2004.
	case mode&os.ModeCharDevice != 0:
		return cChar
	case mode&os.ModeDevice != 0:
		return cDev
	case mode&os.ModeNamedPipe != 0:
		return cFifo
	case mode&os.ModeSymlink != 0:
		return cLink

	// Other file types (though not letters) standardized by POSIX.
	case mode&os.ModeSocket != 0:
		return cSock
	}
	return "?"
}

var buf bytes.Buffer

/* Like filemodestring, but rely only on MODE.  */
func strmode(mode os.FileMode) string {
	buf.Reset()
	buf.WriteString(typeletter(mode))

	if mode&ModeIRUSR != 0 {
		buf.WriteString(cRead)
	} else {
		buf.WriteString(cNone)
	}

	if mode&ModeIWUSR != 0 {
		buf.WriteString(cWrite)
	} else {
		buf.WriteString(cNone)
	}

	if mode&os.ModeSetuid != 0 {
		if mode&ModeIXUSR != 0 {
			buf.WriteString(cUidExec)
		} else {
			buf.WriteString(cUid)
		}
	} else if mode&ModeIXUSR != 0 {
		buf.WriteString(cExec)
	} else {
		buf.WriteString(cNone)
	}

	if mode&ModeIRGRP != 0 {
		buf.WriteString(cRead)
	} else {
		buf.WriteString(cNone)
	}

	if mode&ModeIWGRP != 0 {
		buf.WriteString(cWrite)
	} else {
		buf.WriteString(cNone)
	}

	if mode&os.ModeSetgid != 0 {
		if mode&ModeIXGRP != 0 {
			buf.WriteString(cUidExec)
		} else {
			buf.WriteString(cUid)
		}
	} else if mode&ModeIXGRP != 0 {
		buf.WriteString(cExec)
	} else {
		buf.WriteString(cNone)
	}

	if mode&ModeIROTH != 0 {
		buf.WriteString(cRead)
	} else {
		buf.WriteString(cNone)
	}

	if mode&ModeIWOTH != 0 {
		buf.WriteString(cWrite)
	} else {
		buf.WriteString(cNone)
	}

	if mode&os.ModeSticky != 0 {
		if mode&ModeIXOTH != 0 {
			buf.WriteString(cRes)
		} else {
			buf.WriteString(cResOther)
		}
	} else if mode&ModeIXOTH != 0 {
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
	var buf [32]byte // Mode is uint32.
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
