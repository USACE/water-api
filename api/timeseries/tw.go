package timeseries

import (
	"time"
)

type TimeWindow struct {
	After  time.Time `json:"after"`
	Before time.Time `json:"before"`
}

// CreateTimeWindow
func CreateTimeWindow(a string, b string) (TimeWindow, error) {
	var tw TimeWindow
	var err error
	// If both are missing, then use current time
	if a == "" && b == "" {
		tw.Before = time.Now()
		tw.After = tw.Before.AddDate(0, 0, -7)
		// If one is missing then assumed the other is not
		// additional if...else if used to determine which is missing
	} else if a == "" || b == "" {
		if _, err := time.Parse(time.RFC3339, a); err != nil {
			if tw.Before, err = time.Parse(time.RFC3339, b); err != nil {
				return tw, err
			}
			tw.After = tw.Before.AddDate(0, 0, -7)
		} else if _, err := time.Parse(time.RFC3339, b); err != nil {
			if tw.After, err = time.Parse(time.RFC3339, a); err != nil {
				return tw, err
			}
			tw.Before = tw.After.AddDate(0, 0, 7)
		}
		// Everything else...both are not missing
	} else {
		if tw.After, err = time.Parse(time.RFC3339, a); err != nil {
			return tw, err
		}
		if tw.Before, err = time.Parse(time.RFC3339, b); err != nil {
			return tw, err
		}
	}
	return tw, nil
}
