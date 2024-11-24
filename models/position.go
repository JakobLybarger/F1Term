package models

import "time"

type Position struct {
	Date         time.Time `json:"date"`
	DriverNumber int       `json:"driver_number"`
	Position     int       `json:"position"`
}
