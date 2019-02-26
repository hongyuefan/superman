package strategy

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"sync"

	gin "github.com/gin-gonic/gin"
	"github.com/hongyuefan/superman/config"
	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/protocol"
	"github.com/hongyuefan/superman/solo/skeleton"
	"github.com/hongyuefan/superman/utils"
)

var MapTimeOut map[protocol.KLineType]int64 = map[protocol.KLineType]int64{protocol.SPIDER_TYPE_KLINE_5MIN: 301, protocol.SPIDER_TYPE_KLINE_15MIN: 901, protocol.SPIDER_TYPE_KLINE_30MIN: 1801, protocol.SPIDER_TYPE_KLINE_HOUR: 3601, protocol.SPIDER_TYPE_KLINE_DAY: 86401}

type StratUpDown struct {
	skl          *skeleton.Skeleton
	lock         sync.RWMutex
	touched      map[protocol.KLineType]int64
	chanClose    chan bool
	mapRate      map[protocol.KLineType]float64
	mapCountFlag map[int64]bool
}

func NewStratUpDown() *StratUpDown {
	return &StratUpDown{
		skl:          skeleton.NewSkeleton(),
		chanClose:    make(chan bool, 0),
		mapRate:      make(map[protocol.KLineType]float64, 0),
		touched:      make(map[protocol.KLineType]int64, 0),
		mapCountFlag: make(map[int64]bool, 0),
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

//func (s *StratUpDown) setTouched(b int64) {
//	s.lock.Lock()
//	defer s.lock.Unlock()
//	s.touched = b
//}

//func (s *StratUpDown) getTouched() int64 {
//	s.lock.RLock()
//	defer s.lock.RUnlock()
//	return s.touched
//}

func (s *StratUpDown) isOK(typ protocol.KLineType) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	now := time.Now().Unix()

	if now-s.touched[typ] > MapTimeOut[typ] {

		s.touched[typ] = now

		return true
	}

	return false
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
			s.touched[protocol.SPIDER_TYPE_KLINE_5MIN] = 0
			break
		case 1:
			s.mapRate[protocol.SPIDER_TYPE_KLINE_15MIN] = nRate
			s.touched[protocol.SPIDER_TYPE_KLINE_15MIN] = 0
			break
		case 2:
			s.mapRate[protocol.SPIDER_TYPE_KLINE_30MIN] = nRate
			s.touched[protocol.SPIDER_TYPE_KLINE_30MIN] = 0
			break
		case 3:
			s.mapRate[protocol.SPIDER_TYPE_KLINE_HOUR] = nRate
			s.touched[protocol.SPIDER_TYPE_KLINE_HOUR] = 0
			break
		case 4:
			s.mapRate[protocol.SPIDER_TYPE_KLINE_DAY] = nRate
			s.touched[protocol.SPIDER_TYPE_KLINE_DAY] = 0
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

func (s *StratUpDown) Handler(c *gin.Context) {
}

func (s *StratUpDown) touchMsg(typ protocol.KLineType) {

	kls, ok := s.skl.GetKline(typ, 3)

	if !ok || len(kls) < 2 {
		return
	}

	rate := (kls[0].Close - kls[1].Close) / kls[1].Close * float64(100)

	if typ == protocol.SPIDER_TYPE_KLINE_15MIN {

		rate0 := (kls[1].Close - kls[2].Close) / kls[2].Close * float64(100)

		if rate0 >= 0.5 {

			s.mapCountFlag[kls[1].Time] = true

			if len(s.mapCountFlag) >= 2 && s.isOK(typ) {

				var params0 []string

				params0 = append(params0, typ.String(), "连续上涨", "0.1", fmt.Sprintf("%2.2f", kls[0].Close), config.T.Symbol)

				mobiles := utils.ParseStringToArry(config.T.Mobile, ",")

				for _, mobile := range mobiles {
					if err := utils.SendMsg(config.T.AppID, config.T.AppKey, "86", mobile, params0, config.T.TplId); err != nil {
						logs.Error("send msg error:", err.Error())
					}
				}
			}
		} else {
			s.mapCountFlag = make(map[int64]bool, 0)
		}
	}

	if math.Abs(rate) >= s.mapRate[typ] {

		var params []string
		var ud string

		if rate > 0 {
			ud = "上涨"
		} else {
			ud = "下跌"
		}

		params = append(params, typ.String(), ud, fmt.Sprintf("%2.2f", math.Abs(rate)), fmt.Sprintf("%2.2f", kls[0].Close), config.T.Symbol)

		mobiles := utils.ParseStringToArry(config.T.Mobile, ",")

		if s.isOK(typ) {
			for _, mobile := range mobiles {
				if err := utils.SendMsg(config.T.AppID, config.T.AppKey, "86", mobile, params, config.T.TplId); err != nil {
					logs.Error("send msg error:", err.Error())
				}
			}
		}

	}

}
