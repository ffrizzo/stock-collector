package util

import "time"

//SubtractHours subtract hours on time
func SubtractHours(value time.Time, hours int) time.Time {
	subtract := 1000 * 1000 * 1000 * 60 * 60 * hours * -1
	newTime := value.Add(time.Duration(subtract))
	return newTime
}
