package main

import (
	"fmt"
	"time"
)

// Time units
const (
	Second = 1
	Minute = 60 * Second
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Week   = 7 * Day
	Month  = 30 * Day
	Year   = 12 * Month
)

var now = time.Now().UnixNano()

func reltime(then int64) string {
	diff := (now - then) / 1e9

	switch {
	case diff <= Second:
		return cSecond + "  <s" + cEnd

	case diff < Minute:
		return cSecond +
			fmt.Sprintf("%3d", diff) +
			"s" + cEnd

	case diff < Hour:
		return cMinute +
			fmt.Sprintf("%3d", diff/Minute) +
			"m" + cEnd

	case diff < Hour*36:
		return cHour +
			fmt.Sprintf("%3d", diff/Hour) +
			"h" + cEnd

	case diff < Month:
		return cDay +
			fmt.Sprintf("%3d", diff/Day) +
			"d" + cEnd

	case diff < Year:
		return cWeek +
			fmt.Sprintf("%3d", diff/Week) +
			"w" + cEnd

	//case diff < Year:
	//	return cMonth +
	//		fmt.Sprintf("%3d", diff/Month) +
	//		"mon" + cEnd

	default:
		return cYear +
			fmt.Sprintf("%3d", diff/Year) +
			"y" + cEnd
	}
}
