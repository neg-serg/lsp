package main

import "os"

type indicator int

const (
	TypeLeft indicator = iota
	TypeRight
	TypeEnd
	TypeReset
	TypeNorm
	TypeFile
	TypeDir
	TypeLink
	TypeFifo
	TypeSock
	TypeBlk
	TypeChr
	TypeMissing
	TypeOrphan
	TypeExec
	TypeDoor
	TypeSetuid
	TypeSetgid
	TypeSticky
	TypeOtherWritable
	TypeStickyOtherWritable
	TypeCap
	TypeMultihardlink
	TypeClrToEol
)

var indicatorNames = [...]string{
	"lc", "rc", "ec", "rs", "no", "fi", "di", "ln", "pi", "so",
	"bd", "cd", "mi", "or", "ex", "do", "su", "sg", "st",
	"ow", "tw", "ca", "mh", "cl",
}

func colorType(mode os.FileMode, linkok bool) indicator {
	var t indicator
	if !linkok && isColored(TypeMissing) {
		t = TypeMissing
	} else {
		if mode.IsRegular() {
			t = TypeFile
			switch {
			case mode&os.ModeSetuid != 0 && isColored(TypeSetuid):
				t = TypeSetuid
			case ((mode&os.ModeSetgid) != 0 && isColored(TypeSetgid)):
				t = TypeSetgid
			//case (isColored (C_CAP) && f->has_capability):
			//  t = C_CAP;
			case ((mode&ModeIXUGO) != 0 && isColored(TypeExec)):
				t = TypeExec
				//case ((1 < f->stat.st_nlink) && isColored (C_MULTIHARDLINK)):
				//  t = C_MULTIHARDLINK;
			}
		} else if mode.IsDir() {
			t = TypeDir
			switch {
			case (mode&os.ModeSticky != 0) && (mode&ModeIWOTH != 0) &&
				isColored(TypeStickyOtherWritable):
				t = TypeStickyOtherWritable
			case ((mode&ModeIWOTH) != 0 && isColored(TypeOtherWritable)):
				t = TypeOtherWritable
			case ((mode&os.ModeSticky) != 0 && isColored(TypeSticky)):
				t = TypeSticky
			}
		} else {
			switch {
			case mode&os.ModeSymlink != 0:
				t = TypeLink
			case mode&os.ModeNamedPipe != 0:
				t = TypeFifo
			case mode&os.ModeSocket != 0:
				t = TypeSock
			case mode%os.ModeCharDevice != 0:
				t = TypeChr
			case mode&os.ModeDevice != 0:
				t = TypeBlk
			default:
				// anything else is classified as orphan
				t = TypeOrphan
			}
		}
	}
	if t == TypeLink && !linkok {
		if isColored(TypeOrphan) {
			t = TypeOrphan
		}
	}
	return t
}
