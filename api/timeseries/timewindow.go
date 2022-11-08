package timeseries

// import (
// 	"time"
// )

// type TimeWindow struct {
// 	After  time.Time
// 	Before time.Time
// }

// func NewTimeWindow(a string, b string) (TimeWindow, error) {
// 	now := time.Now()

// 	// If both are missing, return default time window
// 	if a == b == nil {
// 		return TimeWindow{Before: now, After: now.AddDate(0, 0, -7)}, nil
// 	}

// 	// only after is specified
// 	if after != "" {
// 		after, err := time.Parse(time.RFC3339, a)
// 		if err != nil {
// 			return TimeWindow{}, err
// 		}
// 	}

// 	// before

// 	if _, err :=

// 	else if a == "" || b == "" {
// 		if _, err := time.Parse(time.RFC3339, a); err != nil {
// 			if tw.Before, err = time.Parse(time.RFC3339, b); err != nil {
// 				return tw, err
// 			}
// 			tw.After = tw.Before.AddDate(0, 0, -7)
// 		} else if _, err := time.Parse(time.RFC3339, b); err != nil {
// 			if tw.After, err = time.Parse(time.RFC3339, a); err != nil {
// 				return tw, err
// 			}
// 			tw.Before = tw.After.AddDate(0, 0, 7)
// 		}
// 		// Everything else...both are not missing
// 	} else {
// 		if tw.After, err = time.Parse(time.RFC3339, a); err != nil {
// 			return tw, err
// 		}
// 		if tw.Before, err = time.Parse(time.RFC3339, b); err != nil {
// 			return tw, err
// 		}
// 	}
// }

// // CreateTimeWindow
// func CreateTimeWindow(a string, b string) (TimeWindow, error) {

// 	return tw, nil
// }
