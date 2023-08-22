package consumer

type Option func(consumer *Consumer)

func WithOffset(offset int64) Option {
	return func(consumer *Consumer) {
		consumer.Offset = offset
	}
}

func WithPartition(partition int) Option {
	return func(consumer *Consumer) {
		consumer.Partition = partition
	}
}

func WithMDCallback(cb MDCallback) Option {
	return func(consumer *Consumer) {
		consumer.MDCallback = cb
	}
}

func WithSnapshotCallback(cb SnapshotCallback) Option {
	return func(consumer *Consumer) {
		consumer.SnapshotCallback = cb
	}
}
func WithOrderCallback(cb OrderCallback) Option {
	return func(consumer *Consumer) {
		consumer.OrderCallback = cb
	}
}

func WithTransactionCallback(cb TransactionCallback) Option {
	return func(consumer *Consumer) {
		consumer.TransactionCallback = cb
	}
}

func WithChanSize(chanSize int64) Option {
	return func(consumer *Consumer) {
		consumer.ChanSize = chanSize
	}
}
