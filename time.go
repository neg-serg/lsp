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
	const f = "%3d"
	diff := (now - then) / 1e9

	switch {
	case diff <= Second:
		return cSecond + "  <s" + cEnd

	case diff < Minute:
		return cSecond +
			fmt.Sprintf(f, diff) +
			"s" + cEnd

	case diff < Hour:
		return cMinute +
			fmt.Sprintf(f, diff/Minute) +
			"m" + cEnd

	case diff < Hour*36:
		return cHour +
			fmt.Sprintf(f, diff/Hour) +
			"h" + cEnd

	case diff < Month:
		return cDay +
			fmt.Sprintf(f, diff/Day) +
			"d" + cEnd

	case diff < Year:
		return cWeek +
			fmt.Sprintf(f, diff/Week) +
			"w" + cEnd

	//case diff < Year:
	//	return cMonth +
	//		fmt.Sprintf(f, diff/Month) +
	//		"mon" + cEnd

	default:
		return cYear +
			fmt.Sprintf(f, diff/Year) +
			"y" + cEnd
	}
}
