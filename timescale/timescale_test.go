package timescale

import (
	"fmt"
	"testing"
)

func TestNewTimeScale(t *testing.T) {
	timescale := NewTimeScale("09:30:00", "15:00:00", 60)
	fmt.Println(timescale)
}

func TestIntTime2Time(t *testing.T) {
	//timeInt := 93000
	timeInt := 93000111
	t.Logf("%s", IntTime2Time(timeInt))
}
