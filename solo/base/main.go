package base

import (
	"encoding/json"
	"fmt"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/hongyuefan/superman/config"
	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/protocol"
	"github.com/hongyuefan/superman/utils"
	"github.com/okcoin-okex/open-api-v3-sdk/okex-go-sdk-api"
)

type Notice struct {
	Symbol string
	Notice protocol.NoticeType
}

type BaseData struct {
	KLine   *KLineHandler
	Orders  *Order
	Wallets *Wallet
	TTicker *Ticker

	ChanNotice chan Notice
}

func NewBaseData() *BaseData {
	return &BaseData{
		KLine:      NewKLineHandler(),
		TTicker:    NewTicker(),
		ChanNotice: make(chan Notice, 1024),
	}
}

func (b *BaseData) Init() {

	client := okex.NewClient(okex.Config{ApiKey: config.T.ApiKey, SecretKey: config.T.ScretKey, Passphrase: config.T.PassPhrase, Endpoint: config.T.EndPoint, TimeoutSecond: 45, I18n: "en_US", IsPrint: false})

	b.Orders = NewOrder(client)

	if err := b.Orders.LoadOrders(); err != nil {
		logs.Error("load orders error :%s", err.Error())
	}

	b.Wallets = NewWallet(client)

	if err := b.Wallets.LoadCurrency(); err != nil {
		logs.Error("load currency error :%s", err.Error())
	}

	logs.Info("basedata init ok")
}

func (b *BaseData) Handler(js *simplejson.Json) {

	array, err := js.Array()

	if err != nil {
		return
	}

	length := len(array)

	for index := 0; index < length; index++ {

		subJs := js.GetIndex(index)

		channel := subJs.Get("channel").MustString()

		data, ok := subJs.Get("data").CheckGet("errorcode")

		if ok {
			logs.Error("get ws response data error")
			return
		}

		data = subJs.Get("data")

		if err := b.parseDataDetails(channel, data); err != nil {
			logs.Error("parseDataDetails error:%s", err.Error())
			return
		}
	}
	return
}

func (b *BaseData) parseDataDetails(channel string, js *simplejson.Json) error {

	symbol, err := parseSymbol(channel)

	if err != nil {
		return err
	}

	if strings.Contains(channel, "ticker") {

		mticker, err := js.Map()

		if err != nil {
			return err
		}
		timestamp, _ := mticker["timestamp"].(json.Number).Int64()
		high, _ := mticker["high"].(json.Number).Float64()
		vol, _ := mticker["vol"].(json.Number).Float64()
		last, _ := mticker["last"].(json.Number).Float64()
		low, _ := mticker["low"].(json.Number).Float64()
		buy, _ := mticker["buy"].(json.Number).Float64()
		change, _ := mticker["change"].(json.Number).Float64()
		sell, _ := mticker["sell"].(json.Number).Float64()
		dayLow, _ := mticker["dayLow"].(json.Number).Float64()
		dayHigh, _ := mticker["dayHigh"].(json.Number).Float64()

		b.TTicker.SetTicker(high, vol, last, low, buy, change, sell, dayLow, dayHigh, timestamp)

		b.PutNotice(symbol, protocol.NOTICE_TICKER)

	} else if strings.Contains(channel, "depth") {

		typ := getdpkind(channel)

		b.PutNotice(symbol, typ.NoticeType())

	} else if strings.Contains(channel, "kline") {

		typ := getklkind(channel)

		if typ == 0 {
			return fmt.Errorf("kline not suppose %s", channel)
		}

		subArray, _ := js.Array()

		subLen := len(subArray)

		for index := 0; index < subLen; index++ {

			dataJs := js.GetIndex(index)

			v, _ := dataJs.Array()

			b.KLine.Handler(typ, symbol, v[0].(string), v[1].(string), v[2].(string), v[3].(string), v[4].(string), v[5].(string))
		}

		b.PutNotice(symbol, typ.NoticeType())

	}

	return nil
}

func (b *BaseData) PutNotice(symbol string, notice protocol.NoticeType) {
	select {
	case b.ChanNotice <- Notice{Notice: notice, Symbol: symbol}:
	default:
		logs.Error("BaseData Notice Chan Has Full!!")
	}
}

func parseSymbol(channel string) (string, error) {

	strArry := utils.ParseStringToArry(channel, "_")

	if len(strArry) < 5 {
		return "", fmt.Errorf("channel not right %s", channel)
	}
	return strArry[3] + "_" + strArry[4], nil
}

func getdpkind(ch string) protocol.DepthType {

	m := map[string]protocol.DepthType{
		"depth_5": protocol.SPIDER_TYPE_DEPTH_5,
	}

	for k, v := range m {
		if strings.Contains(ch, k) {
			return v
		}
	}

	return 0
}

func getklkind(ch string) protocol.KLineType {

	m := map[string]protocol.KLineType{
		"1min":  protocol.SPIDER_TYPE_KLINE_1MIN,
		"5min":  protocol.SPIDER_TYPE_KLINE_5MIN,
		"15min": protocol.SPIDER_TYPE_KLINE_15MIN,
		"30min": protocol.SPIDER_TYPE_KLINE_30MIN,
		"1hour": protocol.SPIDER_TYPE_KLINE_HOUR,
		"_day":  protocol.SPIDER_TYPE_KLINE_DAY,
		"week":  protocol.SPIDER_TYPE_KLINE_WEEK,
	}
	for k, v := range m {
		if strings.Contains(ch, k) {
			return v
		}
	}
	return 0
}
