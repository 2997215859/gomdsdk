package timescale

import (
	"github.com/golang-module/carbon"
	"sort"
)

// 左闭右闭
type TimeScale struct {
	StartTime  int64
	EndTime    int64
	Scale      int64 // 单位 s
	Minutes    []string
	MinuteSize int
}

func NewTimeScale(startTime, endTime string, scale int) *TimeScale {
	date := carbon.Now().ToDateString()
	startTime = date + " " + startTime
	endTime = date + " " + endTime
	zeroTime := date + " 00:00:00"
	zero := carbon.Parse(zeroTime)

	timescale := &TimeScale{
		StartTime: zero.DiffInSeconds(carbon.Parse(startTime)),
		EndTime:   zero.DiffInSeconds(carbon.Parse(endTime)),
		Scale:     int64(scale), // 单位 s
	}

	for t := timescale.StartTime; t <= timescale.EndTime; t += timescale.Scale {
		timestr := zero.AddSeconds(int(t)).ToTimeString()
		if timestr >= "11:30:00" && timestr < "13:00:00" {
			continue
		}

		timescale.Minutes = append(timescale.Minutes, timestr)
	}
	timescale.MinuteSize = len(timescale.Minutes)
	return timescale
}

func (timescale *TimeScale) GetTi(timestr string) int {
	if len(timestr) != 8 {
		return -1
	}
	idx := sort.Search(timescale.MinuteSize, func(i int) bool { return timescale.Minutes[i] > timestr })
	if idx == timescale.MinuteSize {
		return -1
	}
	return idx
}

func (timescale *TimeScale) Ti2Time(ti int) string {
	ti--
	if ti < 0 || ti >= timescale.MinuteSize {
		return ""
	}
	return timescale.Minutes[ti]
}
