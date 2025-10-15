package lib

import (
	"fmt"
	"strings"
	"time"
)

// Period - Hold From and To timestamps for a given period
type Period struct {
	From time.Time
	To   time.Time
	name string
}

var validPeriods = map[string]func(time.Time) Period{
	"Today":     GetTodayPeriod,
	"ThisMonth": GetThisMonthPeriod,
	"LastMonth": GetLastMonthPeriod,
	"ThisYear":  GetThisYearPeriod,
	"LastYear":  GetLastYearPeriod,
}

func GetValidPeriods() []string {
	keys := make([]string, 0, len(validPeriods))
	for k := range validPeriods {
		keys = append(keys, k)
	}

	return keys
}

func (p *Period) String() string {
	if p.name != "" {
		return p.name
	}
	return p.From.String() + " to " + p.To.String()
}

func (e *Period) Set(v string) error {
	f, ok := validPeriods[v]
	if !ok {
		keys := make([]string, 0, len(validPeriods))
		for k := range validPeriods {
			keys = append(keys, k)
		}
		return fmt.Errorf(`must be one of "%s"`, strings.Join(keys, `", "`))
	}

	*e = f(time.Now().UTC())
	return nil
}

func (e *Period) Type() string {
	return "period"
}

// GetLastMonthPeriod - Get Period for "ThisMonth"
func GetThisMonthPeriod(now time.Time) Period {
	return Period{
		From: time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC),
		To:   time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC),
		name: "ThisMonth",
	}
}

// GetLastMonthPeriod - Get Period for "LastMonth"
func GetLastMonthPeriod(now time.Time) Period {
	firstOfTheMonth := time.Date(now.Year(), now.Month(), 1, 23, 59, 59, 0, time.UTC)

	toTime := firstOfTheMonth.AddDate(0, 0, -1)
	fromTime := time.Date(toTime.Year(), toTime.Month(), 1, 0, 0, 0, 0, time.UTC)

	return Period{
		From: fromTime,
		To:   toTime,
		name: "LastMonth",
	}
}

// GetTodayPeriod - Get Period for "Today"
func GetTodayPeriod(now time.Time) Period {
	return Period{
		From: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
		To:   time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC),
		name: "Today",
	}
}

// GetThisYearPeriod - Get Period for "ThisYear"
func GetThisYearPeriod(now time.Time) Period {
	return Period{
		From: time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC),
		To:   time.Date(now.Year(), 12, 31, 23, 59, 59, 0, time.UTC),
		name: "ThisYear",
	}
}

// GetLastYearPeriod - Get Period for "LastYear"
func GetLastYearPeriod(now time.Time) Period {
	return Period{
		From: time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, time.UTC),
		To:   time.Date(now.Year()-1, 12, 31, 23, 59, 59, 0, time.UTC),
		name: "LastYear",
	}
}
