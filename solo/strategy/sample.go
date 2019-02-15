package strategy

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	gin "github.com/gin-gonic/gin"
	"github.com/hongyuefan/superman/config"
	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/protocol"
	"github.com/hongyuefan/superman/solo/skeleton"
	"github.com/hongyuefan/superman/utils"
	"github.com/okcoin-okex/open-api-v3-sdk/okex-go-sdk-api"
)

type SampleSt struct {
	skl       *skeleton.Skeleton
	client    *okex.Client
	chanClose chan bool
	mapRate   map[protocol.KLineType]float64

	plock     sync.RWMutex
	buyPrice  float64
	sellPrice float64

	slock   sync.RWMutex
	usdtSig bool
	ethSig  bool

	usdt string
	eth  string
}

func NewSampleSt() *SampleSt {
	return &SampleSt{
		skl:       skeleton.NewSkeleton(),
		chanClose: make(chan bool, 0),
		mapRate:   make(map[protocol.KLineType]float64, 0),
		client:    NewOKClient(),
		buyPrice:  0,
		sellPrice: 10000000,
	}
}

func (s *SampleSt) Init() {

	s.skl.Init()

	return
}

func (s *SampleSt) setUsdtSig(sig bool) {
	s.slock.Lock()
	defer s.slock.Unlock()
	s.usdtSig = sig
}

func (s *SampleSt) getUsdtSig() bool {
	s.slock.RLock()
	defer s.slock.RUnlock()
	return s.usdtSig
}

func (s *SampleSt) setEthSig(sig bool) {
	s.slock.Lock()
	defer s.slock.Unlock()
	s.ethSig = sig
}

func (s *SampleSt) getEthSig() bool {
	s.slock.RLock()
	defer s.slock.RUnlock()
	return s.ethSig
}

func (s *SampleSt) getMoneyTimer() {

	timer := time.NewTicker(time.Second * 60)

	for {
		select {
		case <-timer.C:
			if s.getUSDT() >= 1.0 {
				s.setUsdtSig(true)
			} else {
				s.setUsdtSig(false)
			}

			if s.getETH() >= 0.01 {
				s.setEthSig(true)
			} else {
				s.setEthSig(false)
			}

			fmt.Println("usdt balance :", s.usdt, "eth balance :", s.eth)

		case <-s.chanClose:
			timer.Stop()
			return
		}
	}

}

func (s *SampleSt) OnTicker() {

	go s.skl.Run()
	go s.getMoneyTimer()

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

func NewOKClient() *okex.Client {

	var okcf okex.Config

	okcf.Endpoint = config.T.EndPoint
	okcf.ApiKey = config.T.ApiKey
	okcf.SecretKey = config.T.ScretKey
	okcf.Passphrase = config.T.PassPhrase
	okcf.TimeoutSecond = 45
	okcf.IsPrint = false
	okcf.I18n = okex.ENGLISH

	return okex.NewClient(okcf)
}

func (s *SampleSt) setPrice(b, ss float64) {
	s.plock.Lock()
	defer s.plock.Unlock()
	s.buyPrice = b
	s.sellPrice = ss
}

func (s *SampleSt) getPrice() (float64, float64) {
	s.plock.RLock()
	defer s.plock.RUnlock()
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

func (s *SampleSt) getUSDT() float64 {

	var (
		count   int
		balance float64
	)

	for {
		count++
		result, err := s.client.SpotGetAccountCurrency(okex.SpotAccountParams{Currency: "USDT"})
		if err != nil {
			if count > 3 {
				break
			} else {
				continue
			}
		}

		balance, err = strconv.ParseFloat(result.Balance, 64)
		if err != nil {
			if count > 3 {
				break
			} else {
				continue
			}
		}

		s.usdt = result.Balance

		break
	}

	return balance

}
func (s *SampleSt) getETH() float64 {

	var (
		count   int
		balance float64
	)

	for {
		count++
		result, err := s.client.SpotGetAccountCurrency(okex.SpotAccountParams{Currency: "ETH"})
		if err != nil {
			if count > 3 {
				break
			} else {
				continue
			}
		}

		balance, err = strconv.ParseFloat(result.Balance, 64)
		if err != nil {
			if count > 3 {
				break
			} else {
				continue
			}
		}

		s.eth = result.Balance

		break
	}

	return balance
}

func (s *SampleSt) doOrder(tpe, side, instrumentId, size, notional string) bool {

	result, err := s.client.SpotDoOrder(okex.SpotOrderParams{Type: tpe, Side: side, InstrumentId: instrumentId, Size: size, Notional: notional})

	if err != nil {
		fmt.Println("doOrder error:", tpe, side, instrumentId, size, notional, err.Error())
		return false
	}

	return result.Result
}

func (s *SampleSt) touchMsg(typ protocol.KLineType) {

	var params []string

	kls, ok := s.skl.GetKline(typ, 1)

	if !ok || len(kls) < 1 {
		return
	}

	if kls[0].Close <= s.buyPrice && s.getUsdtSig() {

		fmt.Println("touch buy:", kls[0].Close, "usdt sig:", s.getUsdtSig(), "eth sig:", s.getEthSig())

		params = append(params, fmt.Sprintf("%2.2f", s.buyPrice), fmt.Sprintf("%2.2f", kls[0].Close), config.T.Symbol)

		s.sendMsg(278087, params)

		s.doOrder("market", "buy", "eth-usdt", "", s.usdt)

		s.setUsdtSig(false)
	}

	if kls[0].Close >= s.sellPrice && s.getEthSig() {

		fmt.Println("touch sell:", kls[0].Close, "usdt sig:", s.getUsdtSig(), "eth sig:", s.getEthSig())

		rate := (kls[0].Close - s.buyPrice) / s.buyPrice * float64(100)

		params = append(params, fmt.Sprintf("%2.2f", s.sellPrice), fmt.Sprintf("%2.2f", kls[0].Close), fmt.Sprintf("%2.2f", rate), config.T.Symbol)

		s.sendMsg(278086, params)

		s.doOrder("market", "sell", "eth-usdt", s.eth, "")

		s.setEthSig(false)
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
