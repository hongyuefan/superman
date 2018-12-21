package protocol

const TM_LAYOUT_STR = "2006-01-02 15:04:05"

type TickerType int
type DepthType int
type KLineType int

const (
	SPIDER_TYPE_TICKER TickerType = iota + 100
)

const (
	SPIDER_TYPE_DEPTH_5 DepthType = iota + 200
)

const (
	SPIDER_TYPE_KLINE_1MIN KLineType = iota + 300
	SPIDER_TYPE_KLINE_5MIN
	SPIDER_TYPE_KLINE_15MIN
	SPIDER_TYPE_KLINE_30MIN
	SPIDER_TYPE_KLINE_HOUR
	SPIDER_TYPE_KLINE_DAY
	SPIDER_TYPE_KLINE_WEEK
)

const (
	CMD_EXIST             = -1
	CMD_QRY_ACCOUNTS      = 1
	CMD_QRY_ACCOUNT       = 2
	CMD_DO_ORDER          = 3
	CMD_CANCEL_ORDER      = 4
	CMD_QRY_PENDING_ORDER = 5
	CMD_QRY_ORDER         = 6
)

const (
	TOPIC_OKEX_SPIDER_DATA = "okex_spider_data"
	TOPIC_OKEX_ARCHER_REQ  = "okex_archer_req"
	TOPIC_OKEX_ARCHER_RSP  = "okex_archer_rsp"
)

const (
	KL1Min  int32 = 1
	KL3Min  int32 = 2
	KL5Min  int32 = 3
	KL15Min int32 = 4
	KL30Min int32 = 5
	KL1H    int32 = 6
	KL1D    int32 = 7
)
