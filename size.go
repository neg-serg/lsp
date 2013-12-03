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

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

const cSize = ESC + "38;5;216m"

var cSizes = [...]string{
	ESC + "38;5;7;1m" + "B" + cEnd,
	ESC + "38;5;2;1m" + "K" + cEnd,
	ESC + "38;5;14;1m" + "M" + cEnd,
	ESC + "38;5;12;1m" + "G" + cEnd,
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
	e := math.Floor(logn(float64(s), base))
	suffix := cSizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%4.0f"
	if val < 10 {
		f = "%4.1f"
	}

	return fmt.Sprintf("%s"+f+"%s", cSize, val, suffix)
}
