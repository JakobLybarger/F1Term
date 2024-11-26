package models

import (
	"time"
)

type Interval struct {
	Date         time.Time   `json:"date"`
	DriverNumber int         `json:"driver_number"`
	GapToLeader  interface{} `json:"gap_to_leader"`
	Interval     interface{} `json:"interval"`
}
