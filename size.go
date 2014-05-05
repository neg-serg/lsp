package main

import (
	"fmt"
)

func size(w writer, size int64) {
	s := float64(size)
	m := 0
	for s >= 1000.0 {
		m++
		s /= 1024
	}
	w.Write(cSize)
	if s < 10.0 && s != 0 {
		fmt.Fprintf(w, "%3.1f", s)
	} else {
		fmt.Fprintf(w, "%3.0f", s)
	}
	w.Write(cSizes[m])
}
