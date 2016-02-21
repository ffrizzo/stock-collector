package util

import "time"

//SubtractHours subtract X hours on time
func SubtractHours(value time.Time, hours int) time.Time {
	subtract := 1000 * 1000 * 1000 * 60 * 60 * hours * -1
	newTime := value.Add(time.Duration(subtract))
	return newTime
}

//GetStartAndEndTimeOfYesterday return the start/end time for yesterday
func GetStartAndEndTimeOfYesterday() (startTime time.Time, endTime time.Time) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
	end := time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 0, now.Location())
	return start, end
}
