package timescale

import (
	"fmt"
	"testing"
)

func TestNewTimeScale(t *testing.T) {
	timescale := NewTimeScale("09:30:00", "15:00:00", 60)
	fmt.Println(timescale)

	timestr := "08:45:00"
	t.Logf("timestr = %s, ti = %d", timestr, timescale.GetTi(timestr))

	timestr = "09:30:00"
	t.Logf("timestr = %s, ti = %d", timestr, timescale.GetTi(timestr))

	timestr = "09:30:30"
	t.Logf("timestr = %s, ti = %d", timestr, timescale.GetTi(timestr))

	ti := 0
	t.Logf("ti = %d, timestr = %s, ", ti, timescale.Ti2Time(ti))

	ti = 1
	t.Logf("ti = %d, timestr = %s, ", ti, timescale.Ti2Time(ti))

	ti = 2
	t.Logf("ti = %d, timestr = %s, ", ti, timescale.Ti2Time(ti))
}

func TestIntTime2Time(t *testing.T) {
	//timeInt := 93000
	timeInt := 93000111
	t.Logf("%s", IntTime2Time(timeInt))
}
