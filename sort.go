package main

type fileList []*fileInfo
type byVer fileList

func (p byVer) Len() int { return len(p) }
func (p byVer) Less(a, b int) bool {
	aF, bF := p[a], p[b]
	aD, bD := aF.isDir(), bF.isDir()
	if aD != bD {
		return aD
	}
	sA, sB := aF.name, bF.name
	return filevercmp(sA, sB) < 0
}
func (p byVer) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
