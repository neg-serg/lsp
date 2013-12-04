package main

import (
	"fmt"
	"time"
)

const (
	cSecond = cESC + "38;5;12m"
	cMinute = cESC + "38;5;9m"
	cHour   = cESC + "38;5;1m"
	cDay    = cESC + "38;5;8m"
	cWeek   = cESC + "38;5;8m"
	cMonth  = cESC + "38;5;16m"
	cYear   = cESC + "38;5;0m"
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

var now = time.Now()

// Time formats a time into a relative string.
// Time(someT) -> "3 weeks ago"
func reltime(then time.Time) string {
	diff := now.Sub(then) / time.Second

	switch {
	case diff <= Second:
		return cSecond + " <  sec" + cEnd

	case diff < Minute:
		return cSecond +
			fmt.Sprintf("%2d", diff) +
			"  sec" + cEnd

	case diff < Hour:
		return cMinute +
			fmt.Sprintf("%2d", diff/Minute) +
			"  min" + cEnd

	case diff < Hour*36:
		return cHour +
			fmt.Sprintf("%2d", diff/Hour) +
			" hour" + cEnd

	case diff < Month:
		return cDay +
			fmt.Sprintf("%2d", diff/Day) +
			"  day" + cEnd

	//case diff < Month:
	//	return cWeek +
	//		fmt.Sprintf("%-2d", diff/Week) +
	//		" week" + cEnd

	case diff < Year:
		return cMonth +
			fmt.Sprintf("%2d", diff/Month) +
			"  mon" + cEnd

	default:
		return cYear +
			fmt.Sprintf("%2d", diff/Year) +
			" year" + cEnd
	}
}
