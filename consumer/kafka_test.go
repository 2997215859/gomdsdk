package consumer

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/2997215859/gomdsdk/datatype"
)

var snapshotCnt = 0
var orderCnt = 0
var transCnt = 0

func TiCallback(ti int) {
	fmt.Printf("ti: %d...\n", ti)
}

func Callback(md *datatype.MD, meta *datatype.Meta) {
	switch md.Type {
	case datatype.TypeSnapshot:
		if snapshotCnt%10000 == 0 {
			snapshot := md.Data.(*datatype.Snapshot)
			b, err := json.Marshal(snapshot)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}
			fmt.Printf("snapshot: %s\n", string(b))
		}
		snapshotCnt++
	case datatype.TypeOrder:
		if orderCnt%100000 == 0 {
			order := md.Data.(*datatype.Order)
			b, err := json.Marshal(order)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}
			fmt.Printf("order: %s\n", string(b))
		}
		orderCnt++
	case datatype.TypeTransaction:
		if transCnt%100000 == 0 {
			transaction := md.Data.(*datatype.Transaction)
			b, err := json.Marshal(transaction)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}
			fmt.Printf("transaction: %s\n", string(b))
		}
		transCnt++
	}
}

func TransactionCB(d *datatype.Transaction, meta *datatype.Meta) {
	if d.StockID == "000005.SZ" {
		b, err := json.Marshal(d)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}
		fmt.Printf("transaction: %s\n", string(b))
	}
}

func TestKafka(t *testing.T) {
	// var Timescale *timescale.TimeScale = timescale.NewTimeScale("09:30:00", "15:01:00", 60)
	// tiMgr := timgr.NewTiMgr(timgr.WithTiSeqCallback(TiCallback), timgr.WithTimescale(*Timescale))

	consumer := NewConsumer(
		"transaction",
		[]string{"192.168.1.14:19092"},
		WithOffset(0),
		// WithMDCallback(Callback),
		// WithSnapshotTiMgr(tiMgr),
		WithTransactionCallback(TransactionCB),
		WithAuth("md_consumer", "md_c_TT4sHThu"),
	)

	if err := consumer.Run(); err != nil {
		t.Error(err)
		return
	}
}
