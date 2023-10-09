package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/2997215859/gomdsdk/datatype"
	"github.com/2997215859/gomdsdk/timescale"
	"github.com/2997215859/gomdsdk/timgr"
	"github.com/segmentio/kafka-go"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

const DefaultChanSize = 10000000

type SnapshotCallback func(d *datatype.Snapshot, meta *datatype.Meta)
type OrderCallback func(d *datatype.Order, meta *datatype.Meta)
type TransactionCallback func(d *datatype.Transaction, meta *datatype.Meta)
type MDCallback func(d *datatype.MD, meta *datatype.Meta)

type Consumer struct {
	Brokers   []string
	Topic     string
	Offset    int64
	Partition int
	MaxBytes  int

	reader *kafkago.Reader

	SnapshotCallback    SnapshotCallback
	OrderCallback       OrderCallback
	TransactionCallback TransactionCallback
	MDCallback          MDCallback

	mdChannel       chan *kafkago.Message
	snapshotChan    chan *kafkago.Message
	orderChan       chan *kafkago.Message
	transactionChan chan *kafkago.Message

	ChanSize int64
	stopChan chan struct{}

	SnapshotTiMgr    *timgr.TiMgr
	OrderTiMgr       *timgr.TiMgr
	TransactionTiMgr *timgr.TiMgr
	MDTiMgr          *timgr.TiMgr

	Username string
	Password string
}

func NewConsumer(topic string, brokers []string, opts ...Option) *Consumer {
	consumer := &Consumer{
		Brokers:  brokers,
		Topic:    topic,
		Offset:   kafkago.LastOffset,
		MaxBytes: 10e6, // 10MB
		ChanSize: DefaultChanSize,
		stopChan: make(chan struct{}),
	}

	for _, o := range opts {
		o(consumer)
	}

	mechanism, err := scram.Mechanism(scram.SHA256, consumer.Username, consumer.Password)
	if err != nil {
		panic(err)
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	consumer.reader = kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:   consumer.Brokers,
		Topic:     consumer.Topic,
		Partition: consumer.Partition,
		MaxBytes:  consumer.MaxBytes,
		Dialer:    dialer,
	})
	consumer.reader.SetOffset(consumer.Offset)
	return consumer
}

func (c *Consumer) ParseSnapshot(m *kafkago.Message) (*datatype.MD, *datatype.Snapshot, *datatype.Meta, error) {
	meta := &datatype.Meta{
		Key:    string(m.Key),
		Offset: m.Offset,
	}

	snapshot := &datatype.Snapshot{}
	md := &datatype.MD{
		Type: datatype.TypeUnknown,
		Data: snapshot,
	}
	if err := json.Unmarshal(m.Value, &md); err != nil {
		return md, nil, meta, fmt.Errorf("snapshot json.Unmarshal(%+v) error: %s", snapshot, err)
	}
	meta.MDTime = snapshot.Time

	if md.Type != datatype.TypeSnapshot {
		return md, nil, meta, fmt.Errorf("data.type != datatype.TypeSnapshot")
	}

	return md, snapshot, meta, nil
}

func (c *Consumer) ParseOrder(m *kafkago.Message) (*datatype.MD, *datatype.Order, *datatype.Meta, error) {
	meta := &datatype.Meta{
		Key:    string(m.Key),
		Offset: m.Offset,
	}

	order := &datatype.Order{}
	md := &datatype.MD{
		Type: datatype.TypeUnknown,
		Data: order,
	}
	if err := json.Unmarshal(m.Value, &md); err != nil {
		return md, nil, meta, fmt.Errorf("order json.Unmarshal(%+v) error: %s", order, err)
	}
	meta.MDTime = order.Time

	if md.Type != datatype.TypeOrder {
		return md, nil, meta, fmt.Errorf("data.type != datatype.TypeOrder")
	}

	return md, order, meta, nil
}

func (c *Consumer) ParseTransaction(m *kafkago.Message) (*datatype.MD, *datatype.Transaction, *datatype.Meta, error) {
	meta := &datatype.Meta{
		Key:    string(m.Key),
		Offset: m.Offset,
	}

	transaction := &datatype.Transaction{}
	md := &datatype.MD{
		Type: datatype.TypeUnknown,
		Data: transaction,
	}

	if err := json.Unmarshal(m.Value, &md); err != nil {
		return md, nil, meta, fmt.Errorf("transaction json.Unmarshal(%+v) error: %s", transaction, err)
	}
	meta.MDTime = transaction.Time

	if md.Type != datatype.TypeTransaction {
		return md, nil, meta, fmt.Errorf("data.type != datatype.TypeTransaction")
	}

	return md, transaction, meta, nil
}

func (c *Consumer) ParseMD(m *kafkago.Message) (*datatype.MD, *datatype.Meta, error) {
	key := string(m.Key)
	var md *datatype.MD
	var meta *datatype.Meta
	var err error

	switch key {
	case datatype.KeySnapshot:
		md, _, meta, err = c.ParseSnapshot(m)
	case datatype.KeyOrder:
		md, _, meta, err = c.ParseOrder(m)
	case datatype.KeyTransaction:
		md, _, meta, err = c.ParseTransaction(m)
	default:
		return nil, nil, fmt.Errorf("offset(%d) message.Key is unknown", m.Offset)
	}

	if err != nil {
		return md, meta, err
	}

	if md == nil {
		return nil, nil, fmt.Errorf("offset(%d) md is nil", m.Offset)
	}

	if md.Type == datatype.TypeUnknown {
		return md, meta, fmt.Errorf("offset(%d) md.Type is datatype.TypeUnknown", m.Offset)
	}

	return md, meta, nil
}

func (c *Consumer) Handle() {
	if c.SnapshotCallback != nil || c.SnapshotTiMgr != nil {
		c.snapshotChan = make(chan *kafkago.Message, c.ChanSize)

		go func() {
			for {
				select {
				case m := <-c.snapshotChan:
					_, snapshot, meta, err := c.ParseSnapshot(m)
					if err != nil {
						fmt.Printf("c.ParseSnapshot error: %s\n", err)
						break
					}

					if c.SnapshotCallback != nil {
						c.SnapshotCallback(snapshot, meta)
					}

					if c.SnapshotTiMgr != nil {
						c.SnapshotTiMgr.Update(timescale.IntTime2Time(snapshot.Time))
					}
				case <-c.stopChan:
					return
				}
			}
		}()
	}
	if c.OrderCallback != nil || c.OrderTiMgr != nil {
		c.orderChan = make(chan *kafkago.Message, c.ChanSize)
		go func() {
			for {
				select {
				case m := <-c.orderChan:
					_, order, meta, err := c.ParseOrder(m)
					if err != nil {
						fmt.Printf("c.ParseOrder error: %s\n", err)
						break
					}
					if c.OrderCallback != nil {
						c.OrderCallback(order, meta)
					}
					if c.OrderTiMgr != nil {
						c.OrderTiMgr.Update(timescale.IntTime2Time(order.Time))
					}
				case <-c.stopChan:
					return
				}
			}
		}()
	}
	if c.TransactionCallback != nil || c.TransactionTiMgr != nil {
		c.transactionChan = make(chan *kafkago.Message, c.ChanSize)
		go func() {
			for {
				select {
				case m := <-c.transactionChan:
					_, transaction, meta, err := c.ParseTransaction(m)
					if err != nil {
						fmt.Printf("c.ParseTransaction error: %s\n", err)
						break
					}
					if c.TransactionCallback != nil {
						c.TransactionCallback(transaction, meta)
					}
					if c.TransactionTiMgr != nil {
						c.TransactionTiMgr.Update(timescale.IntTime2Time(transaction.Time))
					}
				case <-c.stopChan:
					return
				}
			}
		}()
	}
	if c.MDCallback != nil || c.MDTiMgr != nil {
		c.mdChannel = make(chan *kafkago.Message, c.ChanSize)
		go func() {
			for {
				select {
				case m := <-c.mdChannel:
					md, meta, err := c.ParseMD(m)
					if err != nil {
						fmt.Printf("c.ParseMD error: %s\n", err)
						break
					}
					if c.MDCallback != nil {
						c.MDCallback(md, meta)
					}
					if c.MDTiMgr != nil {
						c.MDTiMgr.Update(timescale.IntTime2Time(meta.MDTime))
					}

				case <-c.stopChan:
					return
				}
			}
		}()
	}
}

func (c *Consumer) ReadMessage(ctx context.Context) error {
	r := c.reader
	m, err := r.ReadMessage(ctx)
	if err != nil {
		return err
	}

	key := string(m.Key)
	if key == "" {
		return fmt.Errorf("offset(%d) key is empty", m.Offset)
	}

	switch key {
	case datatype.KeySnapshot:
		if c.SnapshotCallback != nil || c.snapshotChan != nil {
			c.snapshotChan <- &m
		}
	case datatype.KeyOrder:
		if c.OrderCallback != nil || c.orderChan != nil {
			c.orderChan <- &m
		}
	case datatype.KeyTransaction:
		if c.TransactionCallback != nil || c.transactionChan != nil {
			c.transactionChan <- &m
		}
	default:
		return fmt.Errorf("offset(%d) unkown key(%s)", m.Offset, key)
	}

	if c.MDCallback != nil && c.mdChannel != nil {
		c.mdChannel <- &m
	}
	return nil
}

func (c *Consumer) Read() error {
	ctx := context.Background()
	r := c.reader

	defer func() error {
		if err := r.Close(); err != nil {
			fmt.Printf("failed to close reader: %s\n", err)
			return fmt.Errorf("failed to close reader: %s", err)
		}
		return nil
	}()

	for {
		select {
		case <-c.stopChan:
			return nil
		default:
			if err := c.ReadMessage(ctx); err != nil {
				fmt.Printf("c.ReadMessage: %s\n", err)
			}
		}
	}
}

func (c *Consumer) Run() error {
	go c.Handle()
	return c.Read()
}

func (c *Consumer) Start() {
	go c.Run()
}

func (c *Consumer) Stop() {
	close(c.stopChan)
}
