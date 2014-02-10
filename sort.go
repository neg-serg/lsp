package main

import (
	"sort"
)

type fileList []*fileInfo

type sortFunc func(fileList) sort.Interface

// Len is part of sort.Interface.
func (fl fileList) Len() int { return len(fl) }

// Swap is part of sort.Interface.
func (fl fileList) Swap(i, j int) { fl[i], fl[j] = fl[j], fl[i] }

//

type sizeSort struct{ fileList }

func (sf sizeSort) Less(i, j int) bool {
	a, b := sf.fileList[i], sf.fileList[j]

	if o := byIsDir(a, b); o != 0 {
		return o == -1
	}

	s := a.size - b.size
	if s == 0 {
		return filevercmp(a.name, b.name) < 0
	}
	return s < 0
}

func sortBySize(fl fileList) sort.Interface { return sizeSort{fl} }

//

type timeSort struct{ fileList }

func (sf timeSort) Less(i, j int) bool {
	a, b := sf.fileList[i], sf.fileList[j]

	if o := byIsDir(a, b); o != 0 {
		return o == -1
	}

	s := a.time - b.time
	if s == 0 {
		return filevercmp(a.name, b.name) < 0
	}
	return s > 0
}

func sortByTime(fl fileList) sort.Interface { return timeSort{fl} }

//

type verSort struct{ fileList }

func (sf verSort) Less(i, j int) bool {
	a, b := sf.fileList[i], sf.fileList[j]

	if o := byIsDir(a, b); o != 0 {
		return o < 0
	}

	return filevercmp(a.name, b.name) < 0
}

func sortByVer(fl fileList) sort.Interface { return verSort{fl} }

//

func byIsDir(a, b *fileInfo) int {
	ad, bd := a.isDir(), b.isDir()
	if ad != bd {
		if ad {
			return -1
		}
		return 1
	}
	return 0
}
