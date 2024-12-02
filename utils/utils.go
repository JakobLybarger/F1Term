package utils

import (
	"strconv"

	"github.com/JakobLybarger/formula/models"
)

func GetDriver(drivers []models.Driver, driverNumber int) models.Driver {
	for _, driver := range drivers {
		if driver.Number == driverNumber {
			return driver
		}
	}

	return drivers[0]
}

func GetInterval(intervals []models.Interval, driverNumber int) (models.Interval, bool) {
	for _, interval := range intervals {
		if interval.DriverNumber == driverNumber {
			return interval, true
		}
	}

	return models.Interval{}, false
}

func DisplayAsString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v

	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)

	default:
		return ""
	}
}
