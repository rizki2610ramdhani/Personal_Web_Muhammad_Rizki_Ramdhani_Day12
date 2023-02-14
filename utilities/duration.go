package utilities

import (
	"math"
	"strconv"
	"time"
)

// func to get detail duration
func GetDuration(startDate time.Time, endDate time.Time) string {
	dur := endDate.Sub(startDate) //dur for duration
	durHour := dur.Hours()
	durDay := math.Floor(durHour / 24)
	durWeek := math.Floor(durDay / 7)
	durMonth := math.Floor(durDay / 30)
	durYear := math.Floor(durDay / 365)
	var Duration string

	//check dur
	if durYear > 0 {
		Duration = strconv.FormatFloat(durYear, 'f', 0, 64) + " Tahun"
	} else {
		if durMonth > 0 {
			Duration = strconv.FormatFloat(durMonth, 'f', 0, 64) + " Bulan"
		} else {
			if durWeek > 0 {
				Duration = strconv.FormatFloat(durWeek, 'f', 0, 64) + " Minggu"
			} else {
				if durDay > 0 {
					Duration = strconv.FormatFloat(durDay, 'f', 0, 64) + " Hari"
				} else {
					Duration = "Kurang dari satu hari"
				}
			}
		}
	}

	return Duration
}
