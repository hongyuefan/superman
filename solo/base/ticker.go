package base

import (
	"sync"
)

type Ticker struct {
	mTickerDetail map[string]TickerDetail
	lock          sync.RWMutex
}

func NewTicker() *Ticker {
	return &Ticker{
		mTickerDetail: make(map[string]TickerDetail),
	}
}

func (t *Ticker) GetTicker(symbol string) (TickerDetail, bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	detail, ok := t.mTickerDetail[symbol]
	return detail, ok
}

func (t *Ticker) SetTicker(symbol string, tk TickerDetail) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.mTickerDetail[symbol] = tk
}

/*
buy(double): 买一价。 // 获取全部ticker信息 > best_bid  买一价
high(double): 最高价   //和v3相同
last(double): 最新成交价  //  和v3相同
low(double): 最低价   //和v3相同
sell(double): 卖一价   // 获取ticker信息> best_ask卖一价
timestamp(long)：时间戳   //和v3相同
vol(double): 成交量(最近的24小时)  ///base_volume_24h
timestamp系统时间戳
*/
type TickerDetail struct {
	High      string `json:"high"`
	Vol       string `json:"vol"`
	Last      string `json:"last"`
	Low       string `json:"low"`
	Buy       string `json:"buy"`
	Change    string `json:"change"`
	Sell      string `json:"sell"`
	DayLow    string `json:"dayLow"`
	DayHigh   string `json:"dayHigh"`
	TimeStamp int64  `json:"timestamp"`
}
