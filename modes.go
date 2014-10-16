package main

import (
	"syscall"
)

func typeletter(mode fileMode) []byte {
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
	return []byte("?")
}

func typeletterNoColor(mode fileMode) []byte {
	switch mode & syscall.S_IFMT {
	// these are the most common, so test for them first.
	case syscall.S_IFREG:
		return nNone
	case syscall.S_IFDIR:
		return nDir
	// other letters standardized by POSIX 1003.1-2004.
	case syscall.S_IFCHR:
		return nChar
	case syscall.S_IFBLK:
		return nBlock
	case syscall.S_IFIFO:
		return nFifo
	case syscall.S_IFLNK:
		return nLink
	// other file types (though not letters) standardized by POSIX.
	case syscall.S_IFSOCK:
		return nSock
	}
	return []byte("?")
}

// create mode strings
func strmodeNoColor(buf writer, mode fileMode) {
	buf.Write(typeletterNoColor(mode))
	if mode&syscall.S_IRUSR != 0 {
		buf.Write(nRead)
	} else {
		buf.Write(nNone)
	}

	if mode&syscall.S_IWUSR != 0 {
		buf.Write(nWrite)
	} else {
		buf.Write(nNone)
	}

	if mode&syscall.S_ISUID != 0 {
		if mode&syscall.S_IXUSR != 0 {
			buf.Write(nUIDExec)
		} else {
			buf.Write(nUID)
		}
	} else if mode&syscall.S_IXUSR != 0 {
		buf.Write(nExec)
	} else {
		buf.Write(nNone)
	}

	if mode&syscall.S_IRGRP != 0 {
		buf.Write(nRead)
	} else {
		buf.Write(nNone)
	}

	if mode&syscall.S_IWGRP != 0 {
		buf.Write(nWrite)
	} else {
		buf.Write(nNone)
	}

	if mode&syscall.S_ISGID != 0 {
		if mode&syscall.S_IXGRP != 0 {
			buf.Write(nUIDExec)
		} else {
			buf.Write(nUID)
		}
	} else if mode&syscall.S_IXGRP != 0 {
		buf.Write(nExec)
	} else {
		buf.Write(nNone)
	}

	if mode&syscall.S_IROTH != 0 {
		buf.Write(nRead)
	} else {
		buf.Write(nNone)
	}

	if mode&syscall.S_IWOTH != 0 {
		buf.Write(nWrite)
	} else {
		buf.Write(nNone)
	}

	if mode&syscall.S_ISVTX != 0 {
		if mode&syscall.S_IXOTH != 0 {
			buf.Write(nSticky)
		} else {
			buf.Write(nStickyO)
		}
	} else if mode&syscall.S_IXOTH != 0 {
		buf.Write(nExec)
	} else {
		buf.Write(nNone)
	}
}

// create mode strings
func strmode(buf writer, mode fileMode) {
	buf.Write(typeletter(mode))
	if mode&syscall.S_IRUSR != 0 {
		buf.Write(cRead)
	} else {
		buf.Write(cNone)
	}

	if mode&syscall.S_IWUSR != 0 {
		buf.Write(cWrite)
	} else {
		buf.Write(cNone)
	}

	if mode&syscall.S_ISUID != 0 {
		if mode&syscall.S_IXUSR != 0 {
			buf.Write(cUIDExec)
		} else {
			buf.Write(cUID)
		}
	} else if mode&syscall.S_IXUSR != 0 {
		buf.Write(cExec)
	} else {
		buf.Write(cNone)
	}

	if mode&syscall.S_IRGRP != 0 {
		buf.Write(cRead)
	} else {
		buf.Write(cNone)
	}

	if mode&syscall.S_IWGRP != 0 {
		buf.Write(cWrite)
	} else {
		buf.Write(cNone)
	}

	if mode&syscall.S_ISGID != 0 {
		if mode&syscall.S_IXGRP != 0 {
			buf.Write(cUIDExec)
		} else {
			buf.Write(cUID)
		}
	} else if mode&syscall.S_IXGRP != 0 {
		buf.Write(cExec)
	} else {
		buf.Write(cNone)
	}

	if mode&syscall.S_IROTH != 0 {
		buf.Write(cRead)
	} else {
		buf.Write(cNone)
	}

	if mode&syscall.S_IWOTH != 0 {
		buf.Write(cWrite)
	} else {
		buf.Write(cNone)
	}

	if mode&syscall.S_ISVTX != 0 {
		if mode&syscall.S_IXOTH != 0 {
			buf.Write(cSticky)
		} else {
			buf.Write(cStickyO)
		}
	} else if mode&syscall.S_IXOTH != 0 {
		buf.Write(cExec)
	} else {
		buf.Write(cNone)
	}
}
