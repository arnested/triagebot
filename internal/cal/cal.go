package cal

import (
	"time"

	"github.com/rickar/cal/v2"
	"github.com/rickar/cal/v2/dk"
)

func workCalendar() *cal.BusinessCalendar {
	workcal := cal.NewBusinessCalendar()

	workcal.AddHoliday(dk.Holidays...)
	//nolint:exhaustivestruct
	workcal.AddHoliday(&cal.Holiday{
		Month: time.December,
		Day:   24, //nolint:gomnd
		Func:  cal.CalcDayOfMonth,
	})
	//nolint:exhaustivestruct
	workcal.AddHoliday(&cal.Holiday{
		Month: time.December,
		Day:   31, //nolint:gomnd
		Func:  cal.CalcDayOfMonth,
	})

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
	workdays := cal.WorkdaysInRange(lastThursday, now) - 1

	// This is the first workday if zero workdays have passed and
	// today _is_ a workday.
	return ((workdays == 0) && cal.IsWorkday(now))
}
