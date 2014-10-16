package main

import "fmt"

func writeSize(w writer, size int64, pre []byte, sufs [7][]byte) {
	s := float64(size)
	m := 0
	for s >= 1000.0 {
		m++
		s /= 1024
	}
	w.Write(pre)
	if s != 0 && s < 9.95 {
		fmt.Fprintf(w, "%3.1f", s)
	} else {
		fmt.Fprintf(w, "%3.0f", s)
	}
	w.Write(sufs[m])
}

func size(w writer, size int64) {
	writeSize(w, size, cSize, cSizes)
}

func sizeNoColor(w writer, size int64) {
	writeSize(w, size, nSize, nSizes)
}
