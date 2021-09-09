package timeseries

import "time"

type TimeWindow struct {
	After  time.Time `json:"after"`
	Before time.Time `json:"before"`
}
