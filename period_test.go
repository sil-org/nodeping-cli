package nodeping

import (
	"testing"
	"time"
)

func TestGetPeriod(t *testing.T) {
	date := time.Unix(1502980112, 0).UTC() // Aug 17, 2017
	testCases := []struct {
		period string
		want   Period
	}{
		{
			period: "ThisMonth",
			want: Period{
				From: time.Unix(1501545600, 0).UTC(),
				To:   time.Unix(1503014399, 0).UTC(),
			},
		},
		{
			period: "LastMonth",
			want: Period{
				From: time.Unix(1498867200, 0).UTC(),
				To:   time.Unix(1501545599, 0).UTC(),
			},
		},
		{
			period: "Today",
			want: Period{
				From: time.Unix(1502928000, 0).UTC(),
				To:   time.Unix(1503014399, 0).UTC(),
			},
		},
		{
			period: "ThisYear",
			want: Period{
				From: time.Unix(1483228800, 0).UTC(),
				To:   time.Unix(1514764799, 0).UTC(),
			},
		},
		{
			period: "LastYear",
			want: Period{
				From: time.Unix(1451606400, 0).UTC(),
				To:   time.Unix(1483228799, 0).UTC(),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.period, func(t *testing.T) {
			got := validPeriods[tC.period](date)
			if !got.From.Equal(tC.want.From) {
				t.Errorf("Period 'From' time not correct. Expected %s, got %s", tC.want.From, got.From)
			}

			if !got.To.Equal(tC.want.To) {
				t.Errorf("Period 'To' time not correct. Expected %s, got %s", tC.want.To, got.To)
			}
		})
	}
}
