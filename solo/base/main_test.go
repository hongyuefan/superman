package base

import (
	"testing"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/hongyuefan/superman/protocol"
)

const (
	str = `[{
    "channel":"ok_sub_spot_bch_btc_kline_1min",
    "data":[
        ["1490337840000","995.37","996.75","995.36","996.75","9.112"],
        ["1490337840000","995.37","996.75","995.36","996.75","9.112"]
    ]
}]`

	str1 = `[
    {
        "channel": "ok_sub_spot_bch_btc_ticker",
        "data": {
            "high": "10000",
            "vol": "185.03743858",
            "last": "111",
            "low": "0.00000001",
            "buy": "115",
            "change": "101",
            "sell": "115",
            "dayLow": "0.00000001",
            "dayHigh": "10000",
            "timestamp": 1500444626000
        }
    }
]`
)

func TestBenchMain(t *testing.T) {

	base := NewBaseData()

	base.Init()

	js, _ := simplejson.NewJson([]byte(str))

	base.Handler(js)

	d1, ok1 := base.KLine.Get("bch_btc", protocol.SPIDER_TYPE_KLINE_1MIN, 1)
	d2, ok2 := base.KLine.Get("bch_btc", protocol.SPIDER_TYPE_KLINE_1MIN, 2)
	d3, ok3 := base.KLine.Get("bch_btc", protocol.SPIDER_TYPE_KLINE_1MIN, 3)

	t.Log(d1, ok1, d2, ok2, d3, ok3)

	js, _ = simplejson.NewJson([]byte(str1))

	base.Handler(js)

	t.Log(base.TTicker.GetTicker("bch_btc"))
}
