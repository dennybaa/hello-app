package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

type testCase struct {
	Birthday time.Time
	Today    time.Time
	Days     int
}

func (c testCase) String() string {
	return fmt.Sprintf("Date of birth: %s, today: %s, days to birthday: %d)", c.Birthday, c.Today, c.Days)
}

var tests = []testCase{
	dt("1999-03-10", "2015-03-10", 0),   // Happy Birthday!
	dt("1999-08-26", "2019-08-24", 2),   // birthday is coming this year
	dt("1999-09-10", "2019-08-08", 33),  // birthday is coming this year
	dt("1999-03-10", "2014-08-03", 219), // birthday next year
	dt("1999-03-10", "2015-08-03", 220), // birthday next year which is leap
}

// converts test strings to datetime
func dt(birthday, today string, days int) testCase {
	t1, _ := time.Parse("2006-01-02", birthday)
	t2, _ := time.Parse("2006-01-02", today)
	return testCase{Birthday: t1, Today: t2, Days: days}
}

func TestProcessCases(t *testing.T) {
	for _, c := range tests {
		days, err := DaysToBirthday(c.Birthday, c.Today)
		if err != nil || days != c.Days {
			t.Error(c, ", got: "+strconv.Itoa(days))
		}
	}
}

func TestBirthdayNotBorn(t *testing.T) {
	today := time.Date(1999, time.March, 3, 23, 0, 0, 0, time.UTC)
	birthday := time.Date(2015, time.March, 3, 12, 0, 0, 0, time.UTC)

	if _, err := DaysToBirthday(birthday, today); err == nil {
		t.Error("Not born yet :(")
	}
}
