package main

import (
	"sort"
	"syscall"
)

type fileList []*fileInfo

type sortFunc func(fileList) sort.Interface

// Len is part of sort.Interface.
func (fl fileList) Len() int { return len(fl) }

// Swap is part of sort.Interface.
func (fl fileList) Swap(i, j int) { fl[i], fl[j] = fl[j], fl[i] }

func dirsort(a, b *fileInfo) int {
	ad := a.linkmode&syscall.S_IFMT == syscall.S_IFDIR
	if ad != (b.linkmode&syscall.S_IFMT == syscall.S_IFDIR) {
		if ad {
			return -1
		}
		return 1
	}
	ad = a.mode&syscall.S_IFMT == syscall.S_IFDIR
	if ad != (b.mode&syscall.S_IFMT == syscall.S_IFDIR) {
		if ad {
			return -1
		}
		return 1
	}
	return 0
}

// Size

type sizeSort struct{ fileList }

func (sf sizeSort) Less(i, j int) bool {
	a, b := sf.fileList[i], sf.fileList[j]
	if d := dirsort(a, b); d != 0 {
		return d < 0
	}

	s := a.size - b.size
	if s == 0 {
		return filevercmp(&a.name, &b.name) < 0
	}
	return s > 0
}

func sortBySize(fl fileList) sort.Interface { return sizeSort{fl} }

// Time

type timeSort struct{ fileList }

func (sf timeSort) Less(i, j int) bool {
	a, b := sf.fileList[i], sf.fileList[j]
	if d := dirsort(a, b); d != 0 {
		return d < 0
	}
	if s := a.time - b.time; s != 0 {
		return s > 0
	}
	return filevercmp(&a.name, &b.name) < 0
}

func sortByTime(fl fileList) sort.Interface { return timeSort{fl} }

// Version

type verSort struct{ fileList }

func (sf verSort) Less(i, j int) bool {
	a, b := sf.fileList[i], sf.fileList[j]
	if d := dirsort(a, b); d != 0 {
		return d < 0
	}

	return filevercmp(&a.name, &b.name) < 0
}

func sortByVer(fl fileList) sort.Interface { return verSort{fl} }
