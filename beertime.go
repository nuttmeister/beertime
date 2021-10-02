// Package beertime is used to represent time in beertime format.
package beertime

import (
	"time"
)

const (
	// beerConstant is used to determine if beertime happens
	// on odd or even weeks.
	beerConstant = 1

	// beerTimeStartDay determins what day beertime starts on.
	beerTimeStartDay = time.Friday

	// beerTimeStartHour determins what hour of the day beertime starts.
	beerTimeStartHour = 15
)

// Now returns current beertime based on now.
// Returns bool.
func Now(now time.Time) bool {
	return isItBeerTime(now)
}

// Duration returns the duration until next beertime from now.
// Returns time.Duration.
func Duration(now time.Time) time.Duration {
	return durUntilBeerTime(now, isItBeerTime(now))
}

// isItBeerWeek returns true if current week is a beer week.
// Returns bool.
func isItBeerWeek(week int) bool {
	if week%2 == beerConstant {
		return true
	}
	return false
}

// isItBeerTime will check if it's currently beer time or not.
// Returns bool.
func isItBeerTime(now time.Time) bool {
	_, week := now.ISOWeek()

	if isItBeerWeek(week) {
		if now.Weekday() == beerTimeStartDay && now.Hour() >= beerTimeStartHour {
			return true
		}
	}

	return false
}

// durUntilBeerTime will return the duration until next beertime.
// Returns time.Duration.
func durUntilBeerTime(now time.Time, beerTime bool) time.Duration {
	if beerTime {
		return time.Duration(0)
	}

	// Get which week next beertime will be in.
	days := numDaysToBeerTime(now)
	remaining := remainingDurOfDay(now)

	// Add the static hour of beer time start (16:00 Europe/Stockholm).
	// Add number of full days until next beer time.
	// It will not include the current day and beerTimeStartDay,
	// since those should be calculated with a duration instead.
	// Add the reamning time of the current day.
	next := now.Add(time.Duration(beerTimeStartHour) * time.Hour)
	next = next.Add(time.Duration((days)*24) * time.Hour)
	next = next.Add(remaining)

	return next.Sub(now)
}

// remainingDurOfDay returns the duration remaining for the current day.
// Returns time.Duration.
func remainingDurOfDay(now time.Time) time.Duration {
	length := int(time.Duration(24) * time.Hour)
	cur := (now.Hour() * int(time.Hour)) + (now.Minute() * int(time.Minute)) + (now.Second() * int(time.Second)) + now.Nanosecond()

	return time.Duration(length - cur)
}

// numDaysToBeerTime returns the number of full days until next beer time.
// Returns int.
func numDaysToBeerTime(now time.Time) int {
	_, week := now.ISOWeek()

	switch {
	case !isItBeerWeek(week):
		return 7 + numDaysToBeerDay(now)

	case isItBeerWeek(week) && now.Weekday() > beerTimeStartDay:
		return 14 + numDaysToBeerDay(now)
	}

	return numDaysToBeerDay(now)
}

// numDaysToBeerDay returns the number of days until beerTimeStartDay from current day.
// Returns int.
func numDaysToBeerDay(now time.Time) int {
	return int(beerTimeStartDay - now.Weekday() - 1)
}
