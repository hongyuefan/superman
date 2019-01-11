package strategy

import (
	"fmt"
	"math"
	"strconv"

	"sync"

	"github.com/hongyuefan/superman/config"
	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/protocol"

	"github.com/hongyuefan/superman/solo/skeleton"
	"github.com/hongyuefan/superman/utils"
)

type StratUpDown struct {
	skl       *skeleton.Skeleton
	lock      sync.RWMutex
	chanClose chan bool
	mapRate   map[protocol.KLineType]float64
}

func NewStratUpDown() *StratUpDown {
	return &StratUpDown{
		skl:       skeleton.NewSkeleton(),
		chanClose: make(chan bool, 0),
		mapRate:   make(map[protocol.KLineType]float64, 0),
	}
}

func (s *StratUpDown) Init() {

	s.skl.Init()

	s.initRates(config.T.Rates)

	return
}

func (s *StratUpDown) OnTicker() {

	go s.skl.Run()

	for {
		select {
		case notice := <-s.skl.ChanNotice():
			s.dispatchMsg(notice.Symbol, notice.Notice)
			break
		case <-s.chanClose:
			return
		}
	}
}

func (s *StratUpDown) OnClose() {

	s.skl.Close()

	close(s.chanClose)

}

func (s *StratUpDown) initRates(ss string) {

	arrys := utils.ParseStringToArry(ss, ",")

	for index, rate := range arrys {

		nRate, err := strconv.ParseFloat(rate, 64)
		if err != nil {
			continue
		}

		switch index {
		case 0:
			s.mapRate[protocol.SPIDER_TYPE_KLINE_5MIN] = nRate
			break
		case 1:
			s.mapRate[protocol.SPIDER_TYPE_KLINE_15MIN] = nRate
			break
		case 2:
			s.mapRate[protocol.SPIDER_TYPE_KLINE_30MIN] = nRate
			break
		case 3:
			s.mapRate[protocol.SPIDER_TYPE_KLINE_HOUR] = nRate
			break
		case 4:
			s.mapRate[protocol.SPIDER_TYPE_KLINE_DAY] = nRate
			break

		}
	}
}

func (s *StratUpDown) dispatchMsg(symbol string, notice protocol.NoticeType) {

	switch notice {
	case protocol.NOTICE_KLINE_5MIN:
		s.touchMsg(protocol.SPIDER_TYPE_KLINE_5MIN)
		break
	case protocol.NOTICE_KLINE_15MIN:
		s.touchMsg(protocol.SPIDER_TYPE_KLINE_15MIN)
		break
	case protocol.NOTICE_KLINE_30MIN:
		s.touchMsg(protocol.SPIDER_TYPE_KLINE_30MIN)
		break
	case protocol.NOTICE_KLINE_HOUR:
		s.touchMsg(protocol.SPIDER_TYPE_KLINE_HOUR)
		break
	case protocol.NOTICE_KLINE_DAY:
		s.touchMsg(protocol.SPIDER_TYPE_KLINE_DAY)
		break
	}
}

func (s *StratUpDown) touchMsg(typ protocol.KLineType) {

	kls, ok := s.skl.GetKline(typ, 2)

	if !ok || len(kls) < 2 {
		return
	}

	rate := (kls[0].Close - kls[1].Close) / kls[1].Close * float64(100)

	if math.Abs(rate) >= s.mapRate[typ] {

		var params []string
		var ud string

		if rate > 0 {
			ud = "上涨"
		} else {
			ud = "下跌"
		}

		params = append(params, typ.String(), ud, fmt.Sprintf("%v", math.Abs(rate)), fmt.Sprintf("%v", kls[0].Close), config.T.Symbol)

		mobiles := utils.ParseStringToArry(config.T.Mobile, ",")

		for _, mobile := range mobiles {
			if err := utils.SendMsg(config.T.AppID, config.T.AppKey, "86", mobile, params, config.T.TplId); err != nil {
				logs.Error("send msg error:", err.Error())
			}
		}
	}

}
