package main

import (
	"fmt"
	"time"
)

// Time units
const (
	second = 1
	minute = 60 * second
	hour   = 60 * minute
	day    = 24 * hour
	week   = 7 * day
	month  = 30 * day
	year   = 12 * month
)

var now = time.Now().UnixNano()

func reltime(b writer, then int64) {
	const f = "%3d"
	diff := (now - then) / 1e9

	switch {
	case diff <= second:
		b.WriteString(cSecond + "  <s" + cEnd)

	case diff < minute:
		b.WriteString(cSecond)
		fmt.Fprintf(b, f, diff)
		b.WriteString("s" + cEnd)

	case diff < hour:
		b.WriteString(cMinute)
		fmt.Fprintf(b, f, diff/minute)
		b.WriteString("m" + cEnd)

	case diff < hour*36:
		b.WriteString(cHour)
		fmt.Fprintf(b, f, diff/hour)
		b.WriteString("h" + cEnd)

	case diff < month:
		b.WriteString(cDay)
		fmt.Fprintf(b, f, diff/day)
		b.WriteString("d" + cEnd)

	case diff < year:
		b.WriteString(cWeek)
		fmt.Fprintf(b, f, diff/week)
		b.WriteString("w" + cEnd)

	//case diff < Year:
	//	cMonth +
	//		fmt.Sprintf(f, diff/Month) +
	//		"mon" + cEnd

	default:
		b.WriteString(cYear)
		fmt.Fprintf(b, f, diff/year)
		b.WriteString("y" + cEnd)
	}
}
