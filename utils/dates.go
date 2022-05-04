package utils

import (
	"time"
)

func GetDatesInWeek(date time.Time) []time.Time {
	var datesInWeek = make([]time.Time, 0)
	switch date.Weekday() {
	case time.Sunday:
		for d := 0; d < 7; d++ {
			datesInWeek = append(datesInWeek, newDate(date.Year(), date.Day()+d, date.Month()))
		}
	case time.Monday:
		datesInWeek = append(datesInWeek, newDate(date.Year(), date.Day()-1, date.Month()))
		for d := 0; d < 6; d++ {
			datesInWeek = append(datesInWeek, newDate(date.Year(), date.Day()+d, date.Month()))
		}
	case time.Tuesday:
		for d := 2; d > 0; d-- {
			datesInWeek = append(datesInWeek, newDate(date.Year(), date.Day()-d, date.Month()))
		}
		for d := 0; d < 5; d++ {
			datesInWeek = append(datesInWeek, newDate(date.Year(), date.Day()+d, date.Month()))
		}
	case time.Wednesday:
		for d := 3; d > 0; d-- {
			datesInWeek = append(datesInWeek, newDate(date.Year(), date.Day()-d, date.Month()))
		}
		for d := 0; d < 4; d++ {
			datesInWeek = append(datesInWeek, newDate(date.Year(), date.Day()+d, date.Month()))
		}
	case time.Thursday:
		for d := 4; d > 0; d-- {
			datesInWeek = append(datesInWeek, newDate(date.Year(), date.Day()-d, date.Month()))
		}
		for d := 0; d < 3; d++ {
			datesInWeek = append(datesInWeek, newDate(date.Year(), date.Day()+d, date.Month()))
		}
	case time.Friday:
		for d := 5; d > 0; d-- {
			datesInWeek = append(datesInWeek, newDate(date.Year(), date.Day()-d, date.Month()))
		}
		for d := 0; d < 2; d++ {
			datesInWeek = append(datesInWeek, newDate(date.Year(), date.Day()+d, date.Month()))
		}
	case time.Saturday:
		for d := 6; d >= 0; d-- {
			datesInWeek = append(datesInWeek, newDate(date.Year(), date.Day()-d, date.Month()))
		}
	}
	return datesInWeek
}
func newDate(year, day int, month time.Month) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}
