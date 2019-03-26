package main

import (
	"time"

	"github.com/rickar/cal"
)

func isFirsWorkdaySinceSecurityAnnouncements(c *cal.Calendar, now time.Time) bool {
	// Security issues are announced every Wednesday evening. So
	// we'll handle them at the earliest on Thursday. Calculate
	// our latest Thursday.
	since := ((int(now.Weekday()) - int(time.Thursday)) + 7) % 7
	lastThursday := now.AddDate(0, 0, -since)

	// Calculate how many workdays have passed since last
	// Thursday.
	workdays := c.CountWorkdays(lastThursday, now) - 1

	// This is the first workday if zero workdays have passed and
	// today _is_ a workday.
	return ((workdays == 0) && c.IsWorkday(now))
}
