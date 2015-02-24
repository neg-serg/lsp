package main

import "fmt"

func writeSize(w writer, size int64, sufs [7][]byte) {
	s := float64(size)
	m := 0
	for s >= 1000.0 {
		m++
		s /= 1024
	}
	if s != 0 && s < 9.95 {
		fmt.Fprintf(w, "%.1f", s)
	} else {
		fmt.Fprintf(w, "%3.0f", s)
	}
	w.Write(sufs[m])
}

func size(w writer, size int64) {
	w.Write(cSize)
	writeSize(w, size, cSizes)
}

func sizeNoColor(w writer, size int64) {
	w.Write(nSize)
	writeSize(w, size, nSizes)
}
