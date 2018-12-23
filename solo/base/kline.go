package base

import (
	"strconv"
	"sync"

	"github.com/hongyuefan/superman/protocol"
)

const MaxLen = 10

type KLineHandler struct {
	MKData map[string]*KLineTypeData
}

func NewKLineHandler() *KLineHandler {

	kHandle := new(KLineHandler)

	kHandle.MKData = make(map[string]*KLineTypeData)

	return kHandle
}

func (h *KLineHandler) Set(symbol string, typ protocol.KLineType, klDetail KLineDetail) {

	if h.MKData[symbol] == nil {

		ktyp := NewKLineTypeData()

		kd1min := NewKLineData(60)
		kd5min := NewKLineData(300)
		kd15min := NewKLineData(900)
		kd30min := NewKLineData(1800)
		kd1h := NewKLineData(3600)
		kd1d := NewKLineData(86400)

		ktyp.MKLineData[protocol.SPIDER_TYPE_KLINE_1MIN] = kd1min
		ktyp.MKLineData[protocol.SPIDER_TYPE_KLINE_5MIN] = kd5min
		ktyp.MKLineData[protocol.SPIDER_TYPE_KLINE_15MIN] = kd15min
		ktyp.MKLineData[protocol.SPIDER_TYPE_KLINE_30MIN] = kd30min
		ktyp.MKLineData[protocol.SPIDER_TYPE_KLINE_HOUR] = kd1h
		ktyp.MKLineData[protocol.SPIDER_TYPE_KLINE_DAY] = kd1d

		h.MKData[symbol] = ktyp
	}
	h.MKData[symbol].MKLineData[typ].CleanKLineData(klDetail)
}

func (h *KLineHandler) Get(symbol string, typ protocol.KLineType, index int) (KLineDetail, bool) {
	return h.MKData[symbol].MKLineData[typ].GetKLineData(index)
}

func (h *KLineHandler) Handler(typ protocol.KLineType, symbol, time, open, high, low, close, deal string) {

	nTime, _ := strconv.ParseInt(time, 10, 64)
	nOpen, _ := strconv.ParseFloat(open, 10)
	nHigh, _ := strconv.ParseFloat(high, 10)
	nLow, _ := strconv.ParseFloat(low, 10)
	nClose, _ := strconv.ParseFloat(close, 10)
	nDeal, _ := strconv.ParseFloat(deal, 10)

	klDetail := KLineDetail{
		Type:       typ,
		Time:       nTime,
		Open:       nOpen,
		High:       nHigh,
		Low:        nLow,
		Close:      nClose,
		DealAmount: nDeal,
	}

	h.Set(symbol, typ, klDetail)
}

type KLineTypeData struct {
	MKLineData map[protocol.KLineType]*KLineData
}

func NewKLineTypeData() *KLineTypeData {
	return &KLineTypeData{
		MKLineData: make(map[protocol.KLineType]*KLineData),
	}
}

type KLineData struct {
	MKDetail  map[int]KLineDetail
	lock      sync.RWMutex
	KIntervel int64
}

func NewKLineData(intervel int64) *KLineData {
	return &KLineData{
		MKDetail:  make(map[int]KLineDetail),
		KIntervel: intervel,
	}
}

func (k *KLineData) GetKLineData(index int) (KLineDetail, bool) {
	k.lock.RLock()
	defer k.lock.RUnlock()
	data, ok := k.MKDetail[index]
	return data, ok
}

func (k *KLineData) CleanKLineData(kdetail KLineDetail) {
	k.lock.Lock()
	defer k.lock.Unlock()

	if len(k.MKDetail) == 0 {

		k.MKDetail[1] = kdetail

		return
	}

	curDetail, ok := k.MKDetail[1]
	if !ok {
		panic("kline clean error")
	}

	if (kdetail.Time - curDetail.Time) > k.KIntervel {
		k.ReSort(kdetail)
	} else {
		k.UpdateKLineData(1, kdetail)
	}

	return
}

func (k *KLineData) ReSort(kdetail KLineDetail) {

	nIndex := len(k.MKDetail)

	if nIndex < MaxLen {
		for i := nIndex; i > 0; i-- {
			k.MKDetail[i+1] = k.MKDetail[i]
		}

	} else {
		for i := nIndex; i > 1; i-- {
			k.MKDetail[i] = k.MKDetail[i-1]
		}
	}
	k.MKDetail[1] = kdetail

	return
}

func (k *KLineData) UpdateKLineData(index int, kdetail KLineDetail) {
	k.MKDetail[index] = kdetail
}

type KLineDetail struct {
	Type       protocol.KLineType
	Time       int64
	Open       float64
	High       float64
	Low        float64
	Close      float64
	DealAmount float64
}
