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

func colorType(mode fileMode, linkok bool) indicator {
	var t indicator
	if !linkok {
		t = typeMissing
	} else {
		if mode.isRegular() {
			t = typeFile
			switch {
			case mode&syscall.S_ISUID != 0:
				t = typeSetuid
			case mode&syscall.S_ISGID != 0:
				t = typeSetgid
			case mode& // S_IXUGO
				(syscall.S_IXUSR|syscall.S_IXGRP|syscall.S_IXOTH) != 0:
				t = typeExec
			}
		} else if mode.isDir() {
			t = typeDir
			switch {
			case mode&syscall.S_ISVTX != 0 && mode&syscall.S_IWOTH != 0:
				t = typeStickyOtherWritable
			case mode&syscall.S_IWOTH != 0:
				t = typeOtherWritable
			case mode&syscall.S_ISVTX != 0:
				t = typeSticky
			}
		} else {
			switch mode & syscall.S_IFMT {
			case syscall.S_IFLNK:
				t = typeLink
			case syscall.S_IFIFO:
				t = typeFifo
			case syscall.S_IFSOCK:
				t = typeSock
			case syscall.S_IFCHR:
				t = typeChr
			case syscall.S_IFBLK:
				t = typeBlk
			default:
				// anything else is classified as orphan
				t = typeOrphan
			}
		}
	}
	if t == typeLink && !linkok {
		t = typeOrphan
	}
	return t
}
