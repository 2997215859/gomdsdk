package consumer

import (
	"testing"

	"github.com/2997215859/GoMDSDK/datatype"
)

func Callback(md *datatype.MD, meta *datatype.Meta) {
	
}

func TestKafka(t *testing.T) {
	consumer := NewConsumer(
		"md",
		[]string{"183.134.59.154:9092"},
		WithOffset(0),
		WithMDCallback(Callback),
	)

	if err := consumer.Run(); err != nil {
		t.Error(err)
		return
	}
}
