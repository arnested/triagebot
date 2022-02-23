package cal

import (
	"time"

	"github.com/rickar/cal"
)

func workCalendar() *cal.Calendar {
	workcal := cal.NewCalendar()

	cal.AddDanishHolidays(workcal)
	workcal.AddHoliday(
		cal.DKJuleaften,
		cal.DKNytaarsaften,
	)

	return workcal
}

// IsWorkday or not?
func IsWorkday(now time.Time) bool {
	c := workCalendar()

	return c.IsWorkday(now)
}

const weekLength = 7

// IsFirstWorkdaySinceDrupalSecurityAnnouncements or not?
func IsFirstWorkdaySinceDrupalSecurityAnnouncements(now time.Time) bool {
	cal := workCalendar()

	// Drupal Security issues are announced every Wednesday
	// evening. So we'll handle them at the earliest on
	// Thursday. Calculate our latest Thursday.
	since := ((int(now.Weekday()) - int(time.Thursday)) + weekLength) % weekLength
	lastThursday := now.AddDate(0, 0, -since)

	// Calculate how many workdays have passed since last
	// Thursday.
	workdays := cal.CountWorkdays(lastThursday, now) - 1

	// This is the first workday if zero workdays have passed and
	// today _is_ a workday.
	return ((workdays == 0) && cal.IsWorkday(now))
}
