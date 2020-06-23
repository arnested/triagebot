package cal_test

import (
	"testing"
	"time"

	"arnested.dk/go/triagebot/internal/cal"
)

func TestWorkdays(t *testing.T) {
	var workdays = []struct {
		in  time.Time
		out bool
	}{
		// Christmas
		{time.Date(2018, time.December, 26, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2018, time.December, 27, 8, 40, 0, 0, time.UTC), true},
		{time.Date(2018, time.December, 28, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2018, time.December, 29, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2018, time.December, 30, 8, 40, 0, 0, time.UTC), false},

		// Random week
		{time.Date(2019, time.March, 26, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2019, time.March, 27, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2019, time.March, 28, 8, 40, 0, 0, time.UTC), true},

		// Eastern
		{time.Date(2019, time.April, 17, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2019, time.April, 19, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2019, time.April, 19, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2019, time.April, 20, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2019, time.April, 21, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2019, time.April, 22, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2019, time.April, 23, 8, 40, 0, 0, time.UTC), true},

		// May 1st is Wednesday
		{time.Date(2019, time.May, 1, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2019, time.May, 2, 8, 40, 0, 0, time.UTC), true},

		// Store bededag
		{time.Date(2019, time.May, 15, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2019, time.May, 16, 8, 40, 0, 0, time.UTC), true},
		{time.Date(2019, time.May, 17, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2019, time.May, 18, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2019, time.May, 19, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2019, time.May, 20, 8, 40, 0, 0, time.UTC), false},

		// Kristi Himmelfartsdag
		{time.Date(2019, time.May, 29, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2019, time.May, 30, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2019, time.May, 31, 8, 40, 0, 0, time.UTC), true},
		{time.Date(2019, time.June, 1, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2019, time.June, 2, 8, 40, 0, 0, time.UTC), false},
		{time.Date(2019, time.June, 3, 8, 40, 0, 0, time.UTC), false},
	}

	for _, tt := range workdays {
		out := cal.IsFirstWorkdaySinceDrupalSecurityAnnouncements(tt.in)
		if out != tt.out {
			t.Errorf("%v, should be %v, but got %v", tt.in, tt.out, out)
		}
	}
}
