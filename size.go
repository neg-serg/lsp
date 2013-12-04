package main

import (
	"fmt"
	"math"
)

// IEC Sizes.
// kibis of bits
const (
	Byte = 1
	KiB  = Byte * 1024
	MiB  = KiB * 1024
	GiB  = MiB * 1024
	TiB  = GiB * 1024
	PiB  = TiB * 1024
	EiB  = PiB * 1024
)

const cSize = cESC + "38;5;216m"

var cSizes = [...]string{
	cESC + "38;5;7;1m" + "B" + cEnd,
	cESC + "38;5;2;1m" + "K" + cEnd,
	cESC + "38;5;14;1m" + "M" + cEnd,
	cESC + "38;5;12;1m" + "G" + cEnd,
	cEnd + "T",
	cEnd + "P",
	cEnd + "E",
}

// Human readable bytes
func size(s int64) string {
	const base = 1024
	if s < 10 {
		return fmt.Sprintf("%s%4d%s", cSize, s, cSizes[0])
	}
	f := float64(s)
	e := math.Floor(math.Log(f) / math.Log(base))
	suffix := cSizes[int(e)]
	val := f / math.Pow(base, e)
	if val >= 10 {
		return fmt.Sprintf("%s%4v%s", cSize, int(val), suffix)
	}
	return fmt.Sprintf("%s%4v%s", cSize, math.Floor(val*10)/10, suffix)
}
