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
{"type":1,"stock_id":"600987.SH","trading_day":20230822,"time":144644000,"status":"O","prevclose":7.7,"open":7.74,"high":7.8,"low":7.49,"match":7.59,"ask_prices":[7.59,7.6,7.61,7.62,7.63,7.64,7.65,7.66,7.67,7.68],"ask_volumes":[0.2225,2.37,1.98,2.04,0.91,0.92,0.62,1.6,1.27,1.24],"bid_prices":[7.58,7.57,7.56,7.55,7.54,7.53,7.52,7.51,7.5,7.49],"bid_volumes":[1.94,1.3,1.9,0.55,0.39,0.49,1.1,0.4,1.28,0.41],"trades_num":11141,"volume":5154205,"turnover":39160791,"total_ask_volume":1050841,"total_bid_volume":347300,"weighted_avg_ask_price":7.97,"weighted_avg_bid_price":7.268,"iopv":0,"high_limited":8.47,"low_limited":6.93}
*/
type Snapshot struct {
	StockID    string    // Code: 000001.SZ, 000001.SH
	TradingDay int       // 交易日
	Time       int       // 时间(HHMMSSmmm)
	Status     byte      // 状态
	Prevclose  float64   // 前收盘价
	Open       float64   // 开盘价
	High       float64   // 最高价
	Low        float64   // 最低价
	Match      float64   // 最新价
	AskPrices  []float64 // 申卖价
	AskVolumes []float64 // 申卖量
	BidPrices  []float64 // 申买价
	BidVolumes []float64 // 申买量

	TradesNum           int     // 成交笔数
	Volume              int64   // 成交总量
	Turnover            int64   // 成交总金额
	TotalAskVolume      int64   // 委托卖出总量
	TotalBidVolume      int64   // 委托买入总量
	WeightedAvgAskPrice float64 // 加权平均委卖价格
	WeightedAvgBidPrice float64 // 加权平均委买价格

	IOPV        int     //IOPV净值估值
	HighLimited float64 // 涨停价
	LowLimited  float64 // 跌停价
}

/*
{"type":2,"stock_id":"600348.SH","action_day":20230822,"time":144739440,"order":11965751,"price":7.5,"volume":7000.0,"order_kind":"D","function_code":"S","channel":1,"order_ori_no":11103808,"biz_index":18079243}
*/
type Order struct {
	StockID      string  // Code: 000001.SZ, 000001.SH
	ActionDay    int     // 委托日期(YYMMDD)
	Time         int     // 委托时间(HHMMSSmmm)
	Order        int     // 委托号
	Price        float64 // 委托价格
	Volume       float64 // 委托数量
	OrderKind    byte    // 委托类别
	FunctionCode byte    // 委托代码('B','S','C')
	Channel      int     // channel id
	OrderOriNo   int64   // 原始订单号
	BizIndex     int64   // 业务编号
}

/*
{"type":3,"stock_id":"603196.SH","action_day":20230822,"time":144801100,"index":6539906,"price":180300.0,"volume":100,"turnover":1803,"bsflag":"S","order_kind":" ","function_code":" ","ask_order":11924388,"bid_order":11923476,"channel":5,"biz_index":19366772}
*/
type Transaction struct {
	StockID      string  // Code: 000001.SZ, 000001.SH
	ActionDay    int     // 自然日(YYMMDD)
	Time         int     // 委托时间(HHMMSSmmm)
	Index        int     // 成交编号
	Price        float64 // 成交价格
	Volume       int     // 成交数量
	Turnover     int64   // 成交金额
	Bsflag       byte    // 买卖方向(买：'B', 卖：'S', 不明：' ')
	OrderKind    byte    // 成交类别
	FunctionCode byte    // 成交代码
	AskOrder     int     // 叫卖方委托序号
	BidOrder     int     // 叫买方委托序号
	Channel      int     // channel id
	BizIndex     int64   // 业务编号
}
