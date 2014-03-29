package main

import "syscall"

type indicator int

const (
	typeLeft indicator = iota
	typeRight
	typeEnd
	typeReset
	typeNorm
	typeFile
	typeDir
	typeLink
	typeFifo
	typeSock
	typeBlk
	typeChr
	typeMissing
	typeOrphan
	typeExec
	typeDoor
	typeSetuid
	typeSetgid
	typeSticky
	typeOtherWritable
	typeStickyOtherWritable
	typeCap
	typeMultihardlink
	typeClrToEol
)

func colorType(mode fileMode) indicator {
	switch mode & syscall.S_IFMT {
	case syscall.S_IFREG:
		switch {
		case mode&syscall.S_ISUID != 0:
			return typeSetuid
		case mode&syscall.S_ISGID != 0:
			return typeSetgid
		case mode& // S_IXUGO
			(syscall.S_IXUSR|syscall.S_IXGRP|syscall.S_IXOTH) != 0:
			return typeExec
		}
		return typeFile
	case syscall.S_IFDIR:
		switch {
		case mode&syscall.S_ISVTX != 0 && mode&syscall.S_IWOTH != 0:
			return typeStickyOtherWritable
		case mode&syscall.S_IWOTH != 0:
			return typeOtherWritable
		case mode&syscall.S_ISVTX != 0:
			return typeSticky
		}
		return typeDir
	case syscall.S_IFLNK:
		return typeLink
	case syscall.S_IFIFO:
		return typeFifo
	case syscall.S_IFSOCK:
		return typeSock
	case syscall.S_IFCHR:
		return typeChr
	case syscall.S_IFBLK:
		return typeBlk
	default:
		// anything else is classified as orphan
		return typeOrphan
	}
}
