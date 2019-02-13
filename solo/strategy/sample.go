package strategy

import (
	"fmt"
	"strconv"
	"sync"

	gin "github.com/gin-gonic/gin"
	"github.com/hongyuefan/superman/config"
	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/protocol"
	"github.com/hongyuefan/superman/solo/skeleton"
	"github.com/hongyuefan/superman/utils"
)

type SampleSt struct {
	skl *skeleton.Skeleton

	chanClose chan bool
	mapRate   map[protocol.KLineType]float64

	lock      sync.RWMutex
	buyPrice  float64
	sellPrice float64
}

func NewSampleSt() *SampleSt {
	return &SampleSt{
		skl:       skeleton.NewSkeleton(),
		chanClose: make(chan bool, 0),
		mapRate:   make(map[protocol.KLineType]float64, 0),

		buyPrice:  0,
		sellPrice: 10000000,
	}
}

func (s *SampleSt) Init() {

	s.skl.Init()

	return
}

func (s *SampleSt) OnTicker() {

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

func (s *SampleSt) OnClose() {

	s.skl.Close()

	close(s.chanClose)

}

func (s *SampleSt) setPrice(b, ss float64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.buyPrice = b
	s.sellPrice = ss
}

func (s *SampleSt) getPrice() (float64, float64) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.buyPrice, s.sellPrice
}

func (s *SampleSt) dispatchMsg(symbol string, notice protocol.NoticeType) {

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

func (s *SampleSt) touchMsg(typ protocol.KLineType) {

	var params []string

	kls, ok := s.skl.GetKline(typ, 1)

	if !ok || len(kls) < 1 {
		return
	}

	if kls[0].Close <= s.buyPrice {

		params = append(params, fmt.Sprintf("%2.2f", s.buyPrice), fmt.Sprintf("%2.2f", kls[0].Close), config.T.Symbol)

		s.sendMsg(278087, params)
	}

	if kls[0].Close >= s.sellPrice {

		rate := (kls[0].Close - s.buyPrice) / s.buyPrice * float64(100)

		params = append(params, fmt.Sprintf("%2.2f", s.buyPrice), fmt.Sprintf("%2.2f", kls[0].Close), fmt.Sprintf("%2.2f", rate), config.T.Symbol)

		s.sendMsg(278086, params)
	}

}

func (s *SampleSt) sendMsg(tplId int, params []string) {

	mobiles := utils.ParseStringToArry(config.T.Mobile, ",")

	for _, mobile := range mobiles {
		if err := utils.SendMsg(config.T.AppID, config.T.AppKey, "86", mobile, params, tplId); err != nil {
			logs.Error("send msg error:", err.Error())
		}
	}
}

func (s *SampleSt) Handler(c *gin.Context) {

	var (
		err       error
		buyPrice  float64
		sellPrice float64
	)

	sBuy := c.Query("buy")
	sSell := c.Query("sell")

	if buyPrice, err = strconv.ParseFloat(sBuy, 64); err != nil {
		goto errDeal
	}

	if sellPrice, err = strconv.ParseFloat(sSell, 64); err != nil {
		goto errDeal
	}

	s.setPrice(buyPrice, sellPrice)

	fmt.Println("set buy and sell price:", sBuy, sSell)

	HandleSuccessMsg(c, "Handler", "success")
	return
errDeal:
	HandleErrorMsg(c, "Handler", err.Error())
	return
}

func HandleSuccessMsg(c *gin.Context, requestType, msg string) {
	responseWrite(c, true, msg)
}

func HandleErrorMsg(c *gin.Context, requestType string, result string) {
	msg := fmt.Sprintf("type[%s] From [%s] Error [%s] ", requestType, c.Request.RemoteAddr, result)
	responseWrite(c, false, msg)
}
func responseWrite(ctx *gin.Context, isSuccess bool, result string) {
	ctx.JSON(200, gin.H{
		"isSuccess": isSuccess,
		"message":   result,
	})
}
