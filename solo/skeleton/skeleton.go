package skeleton

import (
	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/protocol"
	"github.com/hongyuefan/superman/solo/base"
	"github.com/hongyuefan/superman/solo/exchanges"
	"github.com/okcoin-okex/open-api-v3-sdk/okex-go-sdk-api"
)

type Skeleton struct {
	baseData *base.BaseData

	exchange *exchanges.OkExChange

	exitChan chan bool
}

func NewSkeleton() *Skeleton {

	skt := new(Skeleton)

	skt.exitChan = make(chan bool, 0)

	skt.baseData = base.NewBaseData()

	skt.exchange = exchanges.NewOkExChange(skt.baseData.Handler)

	return skt
}

func (skt *Skeleton) Init() {

	skt.baseData.Init()

	if err := skt.exchange.Init(); err != nil {
		panic(err)
	}
}

func (skt *Skeleton) Run() {

	skt.exchange.Run()

	<-skt.exitChan
}

func (skt *Skeleton) Close() {
	close(skt.exitChan)
}

func (skt *Skeleton) GetKline(typ protocol.KLineType, count int64) ([]base.KLineDetail, bool) {
	return skt.baseData.KLine.Get(typ, 0, count)
}

func (skt *Skeleton) GetCurrencyNames() []string {

	return skt.baseData.Wallets.GetCurrencyNames()
}

func (skt *Skeleton) GetCurrencyByName(name string) (ok bool, available, balance, hold float64) {

	result, err := skt.baseData.Wallets.GetCurrency(name)

	if err != nil {

		logs.Error("GetCurrencyByName %s error %s", name, err.Error())

		return false, 0, 0, 0
	}

	return true, result.Available, result.Balance, result.Hold
}

func (skt *Skeleton) GetTicker() base.TickerDetail {

	return skt.baseData.TTicker.GetTicker()
}

func (skt *Skeleton) GetPendingOrderIds() []base.OrderQuery {

	return skt.baseData.Orders.GetPendingOrderIds()
}

func (skt *Skeleton) GetOrder(symbol, orderId string) (okex.SpotOrderListResult, error) {

	return skt.baseData.Orders.GetOrder(symbol, orderId)
}

func (skt *Skeleton) DoOrder(clId, typ, side, symbol, margin, price, size, notional string) (okex.SpotOrderResult, error) {

	return skt.baseData.Orders.DoOrder(clId, typ, side, symbol, margin, price, size, notional)
}

func (skt *Skeleton) CanselOrder(symbol, clId, orderId string) (okex.SpotOrderResult, error) {

	return skt.baseData.Orders.CanselOrder(symbol, clId, orderId)
}

func (skt *Skeleton) ChanNotice() chan base.Notice {

	return skt.baseData.ChanNotice
}
