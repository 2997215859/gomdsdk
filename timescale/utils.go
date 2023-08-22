package timescale

import "strconv"

// 91503190 => 09:15:03
// 91503 => 09:15:03
func IntTime2Time(timeInt int) string {
	timestr := strconv.Itoa(timeInt)
	if len(timestr) == 8 || len(timestr) == 5 {
		timestr = "0" + timestr
	}

	if len(timestr) == 9 {
		timestr = timestr[0:6]
	}
	if len(timestr) != 6 {
		return ""
	}
	return timestr[0:2] + ":" + timestr[2:4] + ":" + timestr[4:6]
}
