package data

import "time"

func IsShiftTime(now time.Time) bool {
	hour := now.Hour()

	return hour >= 9 && hour < 18
}
