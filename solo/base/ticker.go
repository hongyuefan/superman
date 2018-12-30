package base

import (
	"sync"
)

type Ticker struct {
	mTickerDetail TickerDetail
	lock          sync.RWMutex
}

func NewTicker() *Ticker {
	return &Ticker{}
}

func (t *Ticker) GetTicker() TickerDetail {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.mTickerDetail
}

func (t *Ticker) SetTicker(high, vol, last, low, buy, change, sell, daylow, dayhigh float64, time int64) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.mTickerDetail = TickerDetail{
		High:      high,
		Vol:       vol,
		Last:      last,
		Low:       low,
		Buy:       buy,
		Change:    change,
		Sell:      sell,
		DayHigh:   dayhigh,
		DayLow:    daylow,
		TimeStamp: time,
	}
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
	High      float64 `json:"high"`
	Vol       float64 `json:"vol"`
	Last      float64 `json:"last"`
	Low       float64 `json:"low"`
	Buy       float64 `json:"buy"`
	Change    float64 `json:"change"`
	Sell      float64 `json:"sell"`
	DayLow    float64 `json:"dayLow"`
	DayHigh   float64 `json:"dayHigh"`
	TimeStamp int64   `json:"timestamp"`
}
