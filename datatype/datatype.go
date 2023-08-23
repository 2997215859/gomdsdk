package datatype

const (
	OrderFuncBuy  = "B"
	OrderFuncSell = "S"
	OrderFunKnown = "C"
)

const (
	OrderKindMkt = "1" // 市价
	OrderKindFix = "2" // 限价

	OrderKindUsf = "U" // 本方最优
	OrderKindUcf = "Y" // 对手方最优
	//OrderKindUtp = "2" // 即时成交
)

const (
	TransactionFuncCancel = "C" // 撤单
	TransactionFuncTrans  = "F" // 成交
)

const (
	TransactionBSFlagBuy     = "B"
	TransactionBSFlagSell    = "S"
	TransactionBSFlagUnknown = "N"
)

const (
	KeySnapshot    = "snapshot"
	KeyOrder       = "order"
	KeyTransaction = "transaction"
	KeyMD          = "md"
)

const (
	TypeUnknown     = 0
	TypeSnapshot    = 1
	TypeOrder       = 2
	TypeTransaction = 3
)

type MD struct {
	Type int
	Data interface{}
}

type Meta struct {
	Key    string
	Offset int64
}

/*
{"type":1,"data":{"stock_id":"000001.SZ","trading_day":20230823,"time":92018000,"status":"I","prevclose":11.37,"open":0.0,"high":0.0,"low":0.0,"match":0.0,"ask_prices":[11.37,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0],"ask_volumes":[9.79,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0],"bid_prices":[11.37,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0],"bid_volumes":[9.79,0.9,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0],"trades_num":0,"volume":0,"turnover":0,"total_ask_volume":0,"total_bid_volume":0,"weighted_avg_ask_price":0.0,"weighted_avg_bid_price":0.0,"iopv":0,"high_limited":12.51,"low_limited":10.23}}
*/
type Snapshot struct {
	StockID    string    `json:"stock_id"`    // Code: 000001.SZ, 000001.SH
	TradingDay int       `json:"trading_day"` // 交易日
	Time       int       `json:"time"`        // 时间(HHMMSSmmm)
	Status     string    `json:"status"`      // 状态
	PrevClose  float64   `json:"prevclose"`   // 前收盘价
	Open       float64   `json:"open"`        // 开盘价
	High       float64   `json:"high"`        // 最高价
	Low        float64   `json:"low"`         // 最低价
	Match      float64   `json:"match"`       // 最新价
	AskPrices  []float64 `json:"ask_prices"`  // 申卖价
	AskVolumes []float64 `json:"ask_volumes"` // 申卖量
	BidPrices  []float64 `json:"bid_prices"`  // 申买价
	BidVolumes []float64 `json:"bid_volumes"` // 申买量

	TradesNum           int     `json:"trading_num"`            // 成交笔数
	Volume              int64   `json:"volume"`                 // 成交总量
	Turnover            int64   `json:"turnover"`               // 成交总金额
	TotalAskVolume      int64   `json:"total_ask_volume"`       // 委托卖出总量
	TotalBidVolume      int64   `json:"total_bid_volume"`       // 委托买入总量
	WeightedAvgAskPrice float64 `json:"weighted_avg_ask_price"` // 加权平均委卖价格
	WeightedAvgBidPrice float64 `json:"weighted_avg_bid_price"` // 加权平均委买价格

	IOPV        int     `json:"iopv"`         //IOPV净值估值
	HighLimited float64 `json:"high_limited"` // 涨停价
	LowLimited  float64 `json:"low_limited"`  // 跌停价
}

/*
order: {"type":2,"data":{"stock_id":"001324.SZ","action_day":20230823,"time":92048490,"order":243233,"price":33.3,"volume":200.0,"order_kind":"0","function_code":"B","channel":2013,"order_ori_no":0,"biz_index":0}
*/
type Order struct {
	StockID      string  `json:"stock_id"`      // Code: 000001.SZ, 000001.SH
	ActionDay    int     `json:"action_day"`    // 委托日期(YYMMDD)
	Time         int     `json:"time"`          // 委托时间(HHMMSSmmm)
	Order        int     `json:"order"`         // 委托号
	Price        float64 `json:"price"`         // 委托价格
	Volume       float64 `json:"volume"`        // 委托数量
	OrderKind    string  `json:"order_kind"`    // 委托类别
	FunctionCode string  `json:"function_code"` // 委托代码('B','S','C')
	Channel      int     `json:"channel"`       // channel id
	OrderOriNo   int64   `json:"order_ori_no"`  // 原始订单号
	BizIndex     int64   `json:"biz_index"`     // 业务编号
}

/*
{"type":3,"data":{"stock_id":"002672.SZ","action_day":20230823,"time":93331810,"index":2535299,"price":0.0,"volume":11700,"turnover":0,"bsflag":" ","order_kind":"0","function_code":"C","ask_order":0,"bid_order":2096531,"channel":2014,"biz_index":0}}
*/
type Transaction struct {
	StockID      string  `json:"stock_id"`      // Code: 000001.SZ, 000001.SH
	ActionDay    int     `json:"action_day"`    // 自然日(YYMMDD)
	Time         int     `json:"time"`          // 委托时间(HHMMSSmmm)
	Index        int     `json:"index"`         // 成交编号
	Price        float64 `json:"price"`         // 成交价格
	Volume       int     `json:"volume"`        // 成交数量
	Turnover     int64   `json:"turnover"`      // 成交金额
	Bsflag       string  `json:"bsflag"`        // 买卖方向(买：'B', 卖：'S', 不明：' ')
	OrderKind    string  `json:"order_kind"`    // 成交类别
	FunctionCode string  `json:"function_code"` // 成交代码
	AskOrder     int     `json:"ask_order"`     // 叫卖方委托序号
	BidOrder     int     `json:"bid_order"`     // 叫买方委托序号
	Channel      int     `json:"channel"`       // channel id
	BizIndex     int64   `json:"biz_index"`     // 业务编号
}
