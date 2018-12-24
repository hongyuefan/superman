package skeleton

import (
	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/protocol"
	"github.com/hongyuefan/superman/solo/base"
	"github.com/hongyuefan/superman/solo/exchanges"
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

func (skt *Skeleton) GetKline(symbol string, typ protocol.KLineType, index int) (ok bool, open, high, low, close, deal float64, time int64) {

	kl, ok := skt.baseData.KLine.Get(symbol, typ, index)

	return ok, kl.Open, kl.High, kl.Low, kl.Close, kl.DealAmount, kl.Time
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

func (skt *Skeleton) GetTicker(symbol string) (base.TickerDetail, bool) {

	return skt.baseData.TTicker.GetTicker(symbol)
}

func (skt *Skeleton) GetPendingOrderIds() []base.OrderQuery {
	return skt.baseData.Orders.GetPendingOrderIds()
}

//func (skt *Skeleton) GetOrderByOrderId()
