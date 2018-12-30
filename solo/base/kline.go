package base

import (
	"strconv"

	"github.com/hongyuefan/superman/protocol"
	"github.com/hongyuefan/superman/solo/database"
)

type KLineHandler struct {
}

func NewKLineHandler() *KLineHandler {
	return new(KLineHandler)
}

func (h *KLineHandler) Get(typ protocol.KLineType, offset, count int64) ([]KLineDetail, bool) {

	query := make(map[string]string, 1)

	kls := make([]KLineDetail, 0)

	switch typ {
	case protocol.SPIDER_TYPE_KLINE_5MIN:

		ml, err := database.GetKLine_5Mins(query, []string{}, []string{"id"}, []string{"desc"}, offset, count)
		if err != nil {
			return kls, false
		}
		for _, m := range ml {

			open, _ := strconv.ParseFloat(m.(database.KLine_5Min).Open, 32)
			high, _ := strconv.ParseFloat(m.(database.KLine_5Min).High, 32)
			low, _ := strconv.ParseFloat(m.(database.KLine_5Min).Low, 32)
			close, _ := strconv.ParseFloat(m.(database.KLine_5Min).Close, 32)
			time, _ := strconv.ParseInt(m.(database.KLine_5Min).Time, 10, 64)
			deal, _ := strconv.ParseFloat(m.(database.KLine_5Min).Deal, 32)

			kl := KLineDetail{
				Type:       protocol.SPIDER_TYPE_KLINE_5MIN,
				Open:       open,
				High:       high,
				Low:        low,
				Close:      close,
				Time:       time,
				DealAmount: deal,
			}

			kls = append(kls, kl)
		}
		return kls, true

	case protocol.SPIDER_TYPE_KLINE_15MIN:
		ml, err := database.GetKLine_15Mins(query, []string{}, []string{"id"}, []string{"desc"}, offset, count)
		if err != nil {
			return kls, false
		}
		for _, m := range ml {

			open, _ := strconv.ParseFloat(m.(database.KLine_15Min).Open, 32)
			high, _ := strconv.ParseFloat(m.(database.KLine_15Min).High, 32)
			low, _ := strconv.ParseFloat(m.(database.KLine_15Min).Low, 32)
			close, _ := strconv.ParseFloat(m.(database.KLine_15Min).Close, 32)
			time, _ := strconv.ParseInt(m.(database.KLine_15Min).Time, 10, 64)
			deal, _ := strconv.ParseFloat(m.(database.KLine_15Min).Deal, 32)

			kl := KLineDetail{
				Type:       protocol.SPIDER_TYPE_KLINE_15MIN,
				Open:       open,
				High:       high,
				Low:        low,
				Close:      close,
				Time:       time,
				DealAmount: deal,
			}

			kls = append(kls, kl)
		}
		return kls, true
	case protocol.SPIDER_TYPE_KLINE_30MIN:
		ml, err := database.GetKLine_30Mins(query, []string{}, []string{"id"}, []string{"desc"}, offset, count)
		if err != nil {
			return kls, false
		}
		for _, m := range ml {

			open, _ := strconv.ParseFloat(m.(database.KLine_30Min).Open, 32)
			high, _ := strconv.ParseFloat(m.(database.KLine_30Min).High, 32)
			low, _ := strconv.ParseFloat(m.(database.KLine_30Min).Low, 32)
			close, _ := strconv.ParseFloat(m.(database.KLine_30Min).Close, 32)
			time, _ := strconv.ParseInt(m.(database.KLine_30Min).Time, 10, 64)
			deal, _ := strconv.ParseFloat(m.(database.KLine_30Min).Deal, 32)

			kl := KLineDetail{
				Type:       protocol.SPIDER_TYPE_KLINE_30MIN,
				Open:       open,
				High:       high,
				Low:        low,
				Close:      close,
				Time:       time,
				DealAmount: deal,
			}

			kls = append(kls, kl)
		}
		return kls, true
	case protocol.SPIDER_TYPE_KLINE_HOUR:
		ml, err := database.GetKLine_Hours(query, []string{}, []string{"id"}, []string{"desc"}, offset, count)
		if err != nil {
			return kls, false
		}
		for _, m := range ml {

			open, _ := strconv.ParseFloat(m.(database.KLine_Hour).Open, 32)
			high, _ := strconv.ParseFloat(m.(database.KLine_Hour).High, 32)
			low, _ := strconv.ParseFloat(m.(database.KLine_Hour).Low, 32)
			close, _ := strconv.ParseFloat(m.(database.KLine_Hour).Close, 32)
			time, _ := strconv.ParseInt(m.(database.KLine_Hour).Time, 10, 64)
			deal, _ := strconv.ParseFloat(m.(database.KLine_Hour).Deal, 32)

			kl := KLineDetail{
				Type:       protocol.SPIDER_TYPE_KLINE_HOUR,
				Open:       open,
				High:       high,
				Low:        low,
				Close:      close,
				Time:       time,
				DealAmount: deal,
			}

			kls = append(kls, kl)
		}
		return kls, true
	case protocol.SPIDER_TYPE_KLINE_DAY:
		ml, err := database.GetKLine_Days(query, []string{}, []string{"id"}, []string{"desc"}, offset, count)
		if err != nil {
			return kls, false
		}
		for _, m := range ml {

			open, _ := strconv.ParseFloat(m.(database.KLine_Day).Open, 32)
			high, _ := strconv.ParseFloat(m.(database.KLine_Day).High, 32)
			low, _ := strconv.ParseFloat(m.(database.KLine_Day).Low, 32)
			close, _ := strconv.ParseFloat(m.(database.KLine_Day).Close, 32)
			time, _ := strconv.ParseInt(m.(database.KLine_Day).Time, 10, 64)
			deal, _ := strconv.ParseFloat(m.(database.KLine_Day).Deal, 32)

			kl := KLineDetail{
				Type:       protocol.SPIDER_TYPE_KLINE_DAY,
				Open:       open,
				High:       high,
				Low:        low,
				Close:      close,
				Time:       time,
				DealAmount: deal,
			}

			kls = append(kls, kl)
		}
		return kls, true
	}
	return kls, false
}

func (h *KLineHandler) Handler(typ protocol.KLineType, symbol, time, open, high, low, close, deal string) {

	switch typ {
	case protocol.SPIDER_TYPE_KLINE_5MIN:
		database.SetKLine_5MinByTime(&database.KLine_5Min{Open: open, High: high, Low: low, Close: close, Deal: deal, Time: time})
		break
	case protocol.SPIDER_TYPE_KLINE_15MIN:
		database.SetKLine_15MinByTime(&database.KLine_15Min{Open: open, High: high, Low: low, Close: close, Deal: deal, Time: time})
		break
	case protocol.SPIDER_TYPE_KLINE_30MIN:
		database.SetKLine_30MinByTime(&database.KLine_30Min{Open: open, High: high, Low: low, Close: close, Deal: deal, Time: time})
		break
	case protocol.SPIDER_TYPE_KLINE_HOUR:
		database.SetKLine_HourByTime(&database.KLine_Hour{Open: open, High: high, Low: low, Close: close, Deal: deal, Time: time})
		break
	case protocol.SPIDER_TYPE_KLINE_DAY:
		database.SetKLine_DayByTime(&database.KLine_Day{Open: open, High: high, Low: low, Close: close, Deal: deal, Time: time})
		break
	}

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
