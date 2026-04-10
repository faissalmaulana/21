package service

import "time"

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func TimeToDate(t time.Time) Date {
	y, m, d := t.Date()

	return Date{
		Year:  y,
		Month: m,
		Day:   d,
	}
}
