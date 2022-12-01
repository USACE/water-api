package models

import (
	"time"
)

type TimeWindow struct {
	After  time.Time
	Before time.Time
}

func NewDefaultTimeWindow() TimeWindow {
	now := time.Now()
	return TimeWindow{Before: now, After: now.AddDate(0, 0, -7)}
}

func NewTimeWindow(pTW *TimeWindow, a string, b string) error {

	// If both are missing, return default time window
	if a == "" && b == "" {
		*pTW = NewDefaultTimeWindow()
		return nil
	}

	tw := TimeWindow{}

	// After
	if a != "" {
		after, err := time.Parse(time.RFC3339, a)
		if err != nil {
			return err
		}
		tw.After = after
	}

	// Before
	if b != "" {
		before, err := time.Parse(time.RFC3339, b)
		if err != nil {
			return err
		}
		tw.Before = before
	}

	*pTW = tw

	return nil
}
