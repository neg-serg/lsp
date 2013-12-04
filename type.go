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
	if !linkok && isColored(typeMissing) {
		t = typeMissing
	} else {
		if mode.isRegular() {
			t = typeFile
			switch {
			case mode&syscall.S_ISUID != 0 && isColored(typeSetuid):
				t = typeSetuid
			case ((mode&syscall.S_ISGID) != 0 && isColored(typeSetgid)):
				t = typeSetgid
			//case (isColored (C_CAP) && f->has_capability):
			//  t = C_CAP;
			case ((mode&modeIXUGO) != 0 && isColored(typeExec)):
				t = typeExec
				//case ((1 < f->stat.st_nlink) && isColored (C_MULTIHARDLINK)):
				//  t = C_MULTIHARDLINK;
			}
		} else if mode.isDir() {
			t = typeDir
			switch {
			case (mode&syscall.S_ISVTX != 0) && (mode&modeIWOTH != 0) &&
				isColored(typeStickyOtherWritable):
				t = typeStickyOtherWritable
			case ((mode&modeIWOTH) != 0 && isColored(typeOtherWritable)):
				t = typeOtherWritable
			case ((mode&syscall.S_ISVTX) != 0 && isColored(typeSticky)):
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
		if isColored(typeOrphan) {
			t = typeOrphan
		}
	}
	return t
}
