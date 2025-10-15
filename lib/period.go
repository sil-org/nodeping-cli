package lib

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

const (
	periodThisMonth PeriodValue = "ThisMonth"
	periodLastMonth PeriodValue = "LastMonth"
	periodThisYear  PeriodValue = "ThisYear"
	periodLastYear  PeriodValue = "LastYear"
)

type PeriodValue string

func (e *PeriodValue) String() string {
	return string(*e)
}

func (e *PeriodValue) Set(v string) error {
	if !IsPeriodValid(v) {
		periods := GetValidPeriods()
		valid := make([]string, len(periods))
		for i, v := range periods {
			valid[i] = string(v)
		}
		return fmt.Errorf(`must be one of "%s"`, strings.Join(valid, `", "`))
	}

	*e = PeriodValue(v)
	return nil
}

func (e *PeriodValue) Type() string {
	return "period"
}

func GetValidPeriods() []PeriodValue {
	return []PeriodValue{
		periodThisMonth,
		periodLastMonth,
		periodThisYear,
		periodLastYear,
	}
}

func IsPeriodValid(pType string) bool {
	return slices.Contains(GetValidPeriods(), PeriodValue(pType))
}

// Period - Hold From and To timestamps for a given period
type Period struct {
	From int64
	To   int64
}

func (p *Period) String() (string, string) {
	return time.Unix(int64(p.From), 0).UTC().String(), time.Unix(int64(p.To), 0).UTC().String()
}

// GetLastMonthPeriod - Get Period for "ThisMonth"
func GetThisMonthPeriod(now time.Time) Period {
	fromTime := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	toTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)

	return Period{
		From: fromTime.Unix(),
		To:   toTime.Unix(),
	}
}

// GetLastMonthPeriod - Get Period for "LastMonth"
func GetLastMonthPeriod(now time.Time) Period {
	firstOfTheMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastDayOfLastMonth := firstOfTheMonth.AddDate(0, 0, -1)

	fromTime := time.Date(lastDayOfLastMonth.Year(), lastDayOfLastMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	toTime := time.Date(lastDayOfLastMonth.Year(), lastDayOfLastMonth.Month(), lastDayOfLastMonth.Day(), 23, 59, 59, 0, time.UTC)

	return Period{
		From: fromTime.Unix(),
		To:   toTime.Unix(),
	}
}

// GetTodayPeriod - Get Period for "Today"
func GetTodayPeriod(now time.Time) Period {
	fromTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	toTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)

	return Period{
		From: fromTime.Unix(),
		To:   toTime.Unix(),
	}
}

// GetThisYearPeriod - Get Period for "ThisYear"
func GetThisYearPeriod(now time.Time) Period {
	fromTime := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	toTime := time.Date(now.Year(), 12, 31, 23, 59, 59, 0, time.UTC)

	return Period{
		From: fromTime.Unix(),
		To:   toTime.Unix(),
	}
}

// GetLastYearPeriod - Get Period for "LastYear"
func GetLastYearPeriod(now time.Time) Period {
	fromTime := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, time.UTC)
	toTime := time.Date(now.Year()-1, 12, 31, 23, 59, 59, 0, time.UTC)

	return Period{
		From: fromTime.Unix(),
		To:   toTime.Unix(),
	}
}

// GetPeriodByName - Get Period by name
func GetPeriodByName(name PeriodValue, now int64) Period {
	// Initialize "now" in case it was not provided
	if now == 0 {
		now = time.Now().UTC().Unix()
	}
	date := time.Unix(now, 0)

	switch name {
	case periodThisMonth:
		return GetThisMonthPeriod(date)
	case periodLastMonth:
		return GetLastMonthPeriod(date)
	case periodThisYear:
		return GetThisYearPeriod(date)
	case periodLastYear:
		return GetLastYearPeriod(date)
	default:
		return GetTodayPeriod(date)
	}
}
