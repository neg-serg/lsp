package main

import "os"

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

func colorType(mode os.FileMode, linkok bool) indicator {
	var t indicator
	if !linkok && isColored(typeMissing) {
		t = typeMissing
	} else {
		if mode.IsRegular() {
			t = typeFile
			switch {
			case mode&os.ModeSetuid != 0 && isColored(typeSetuid):
				t = typeSetuid
			case ((mode&os.ModeSetgid) != 0 && isColored(typeSetgid)):
				t = typeSetgid
			//case (isColored (C_CAP) && f->has_capability):
			//  t = C_CAP;
			case ((mode&modeIXUGO) != 0 && isColored(typeExec)):
				t = typeExec
				//case ((1 < f->stat.st_nlink) && isColored (C_MULTIHARDLINK)):
				//  t = C_MULTIHARDLINK;
			}
		} else if mode.IsDir() {
			t = typeDir
			switch {
			case (mode&os.ModeSticky != 0) && (mode&modeIWOTH != 0) &&
				isColored(typeStickyOtherWritable):
				t = typeStickyOtherWritable
			case ((mode&modeIWOTH) != 0 && isColored(typeOtherWritable)):
				t = typeOtherWritable
			case ((mode&os.ModeSticky) != 0 && isColored(typeSticky)):
				t = typeSticky
			}
		} else {
			switch {
			case mode&os.ModeSymlink != 0:
				t = typeLink
			case mode&os.ModeNamedPipe != 0:
				t = typeFifo
			case mode&os.ModeSocket != 0:
				t = typeSock
			case mode%os.ModeCharDevice != 0:
				t = typeChr
			case mode&os.ModeDevice != 0:
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
