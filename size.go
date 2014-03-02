package main

import (
	"fmt"
	"math"
)

// IEC Sizes; kibis of bits
const (
	Byte = 1
	KiB  = Byte * 1024
	MiB  = KiB * 1024
	GiB  = MiB * 1024
	TiB  = GiB * 1024
	PiB  = TiB * 1024
	EiB  = PiB * 1024
)

// Human readable bytes
func size(b writer, s int64) {
	const base = 1024
	if s < 10 {
		fmt.Fprintf(b, "%s%3d%s", cSize, s, cSizes[0])
		return
	}
	f := float64(s)
	e := math.Floor(math.Log(f) / math.Log(base))
	suffix := cSizes[int(e)]
	val := f / math.Pow(base, e)
	if val >= 10 {
		fmt.Fprintf(b, "%s%3v%s", cSize, int(val), suffix)
	} else {
		fmt.Fprintf(b, "%s%3v%s", cSize, math.Floor(val*10)/10, suffix)
	}
}
