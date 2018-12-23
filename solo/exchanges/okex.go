package exchanges

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/models"
	"github.com/hongyuefan/superman/protocol"
	"github.com/hongyuefan/superman/utils"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
)

type OkExChange struct {
	wsUrl         string
	mapErr        map[int]string
	okexSymbols   []string
	okexKlineTime []string
	okexDepth     []string
	hbInterval    int64
}

func NewOkExChange() Exchange {
	return &OkExChange{
		mapErr: make(map[int]string),
	}
}

func (t *OkExChange) Init() error {

	spider, err := models.GetSpiderByName("okex")

	if err != nil {
		logs.Error("Get SqlDB Spider Okex Error:%s", err.Error())
		return err
	}

	t.wsUrl = spider.WsUrl
	t.okexSymbols = utils.ParseStringToArry(spider.Symbols, ",")
	t.okexKlineTime = utils.ParseStringToArry(spider.KlineTime, ",")
	t.okexDepth = utils.ParseStringToArry(spider.Depth, ",")
	t.hbInterval = spider.HeartBeat

	utils.InitOkexErrorMap(t.mapErr)

	return nil
}

func (t *OkExChange) Run() {

	for _, s := range t.okexSymbols {
		go t.RunImpl(s)
	}
}

/*
  主协程开出读写2个协程，并监控他们是否退出，只要有一个退出
  主协程会结束链接，这2个协程遇到链接结束肯定会退出，主协程重新来过
*/
func (t *OkExChange) RunImpl(symbol string) {

	for {

		c := utils.Reconnect(t.wsUrl, "okex", "okex_ws")

		rgc := make(chan int)
		wgc := make(chan int)

		go readLoop(c, rgc, symbol)
		go writeLoop(c, wgc, symbol, t.okexKlineTime, t.okexDepth, t.hbInterval)

	L:
		for {
			select {
			case _, ok := <-rgc:
				if !ok {
					break L
				}
			case _, ok := <-wgc:
				if !ok {
					break L
				}
			}
		}

		c.Close()

		logs.Error("okex run %s_%s restart.... ", symbol)
	}
}

func gzipDecode(in []byte) ([]byte, error) {

	reader := flate.NewReader(bytes.NewReader(in))

	defer reader.Close()

	return ioutil.ReadAll(reader)

}

func readLoop(c *websocket.Conn, rgc chan int, symbol string) {

	defer close(rgc)

	var (
		messageType int
		message     []byte
		err         error
	)

	for {
		if messageType, message, err = c.ReadMessage(); err != nil {
			logs.Error("%s_%s sub ws error read:%s", symbol, err.Error())
			return
		}

		switch messageType {

		case websocket.BinaryMessage:

			message, err = gzipDecode(message)

			if err != nil {
				logs.Error("ws readmessage error: %s", err.Error())
				continue
			}
		}

		fmt.Println("ws receive message: ", string(message))

		// 去除心跳回应
		if len(message) == len(`{"event":"pong"}`) {
			continue
		}

		js, err := simplejson.NewJson(message)

		if err != nil {
			logs.Error("%s_%s sub ws parse json error:%s, json: %s", symbol, err.Error(), message)
			return
		}

		if _, err := ParseReply(js); err != nil {
			logs.Error("%s_%s sub ws parse json error:%s, json: %s", symbol, err.Error(), message)
			return
		}
	}
}

func writeLoop(c *websocket.Conn, wgc chan int, symbol string, kline, depth []string, hb int64) {

	tc := time.NewTimer(time.Duration(hb) * time.Second)
	defer tc.Stop()

	defer close(wgc)

	subSpotTicker(c, symbol)
	subSpotDepth(c, symbol, depth)
	subSpotKline(c, symbol, kline)

	for {

		err := c.WriteMessage(websocket.TextMessage, []byte("{'event':'ping'}"))

		if err != nil {
			logs.Error("okex write goroutine write error, %s", err.Error())
			return
		}

		<-tc.C

		tc.Reset(time.Duration(hb) * time.Second)
	}
}

func subSpotTicker(c *websocket.Conn, symbol string) {

	indexStr := "{'event':'addChannel','channel':'ok_sub_spot_%s_ticker'}"

	symbolIndex := fmt.Sprintf(indexStr, symbol)

	c.WriteMessage(websocket.TextMessage, []byte(symbolIndex))

}

func subSpotDepth(c *websocket.Conn, symbol string, nums []string) {

	depthStr := "{'event':'addChannel','channel':'ok_sub_spot_%s_depth_%s'}"

	for _, num := range nums {

		depth := fmt.Sprintf(depthStr, symbol, num)

		c.WriteMessage(websocket.TextMessage, []byte(depth))
	}
}

func subSpotKline(c *websocket.Conn, symbol string, keys []string) {

	klineStr := "{'event':'addChannel','channel':'ok_sub_spot_%s_kline_%s'}"

	for _, key := range keys {

		kline := fmt.Sprintf(klineStr, symbol, key)

		c.WriteMessage(websocket.TextMessage, []byte(kline))
	}
}

func ParseReply(js *simplejson.Json) (bool, error) {

	array, err := js.Array()

	if err != nil {
		return false, err
	}

	length := len(array)

	for index := 0; index < length; index++ {

		subJs := js.GetIndex(index)

		data, ok := subJs.Get("data").CheckGet("errorcode")
		if ok {
			return false, fmt.Errorf("get ticker data error")
		}

		channel := subJs.Get("channel").MustString()

		data = subJs.Get("data")

		if err := parseDataDetails(channel, data); err != nil {
			return false, err
		}
	}

	return true, nil
}

func parseDataDetails(channel string, js *simplejson.Json) error {

	if strings.Contains(channel, "ticker") {

		okexSpiderData(protocol.SPIDER_TYPE_TICKER, js.MustString())

	}

	if strings.Contains(channel, "depth") {

		typ := getdpkind(channel)

		if typ != 0 {
			okexSpiderData(typ, js.MustString())
		}
	}

	if strings.Contains(channel, "kline") {

		typ := getklkind(channel)

		if typ != 0 {
			okexSpiderData(typ, js.MustString())
		}
	}

	return nil
}

func getdpkind(ch string) int {

	m := map[string]int{
		"depth_5": protocol.SPIDER_TYPE_DEPTH_5,
	}

	for k, v := range m {
		if strings.Contains(ch, k) {
			return v
		}
	}

	return 0
}

func getklkind(ch string) int {

	m := map[string]int{
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

func okexSpiderData(typ int, msg string) error {
	return utils.PackAndReplyToBroker(protocol.TOPIC_OKEX_SPIDER_DATA, "okex", typ, msg)
}
