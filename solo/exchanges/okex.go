package exchanges

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io/ioutil"

	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
	"github.com/hongyuefan/superman/config"
	"github.com/hongyuefan/superman/logs"

	"github.com/hongyuefan/superman/utils"
)

type MsgHandler func(*simplejson.Json)

type OkExChange struct {
	wsUrl         string
	mapErr        map[int]string
	okexSymbols   []string
	okexKlineTime []string
	okexDepth     []string
	hbInterval    int64

	msgHandler MsgHandler
}

func NewOkExChange(msgHandler MsgHandler) *OkExChange {
	return &OkExChange{
		mapErr:     make(map[int]string),
		msgHandler: msgHandler,
	}
}

func (t *OkExChange) Init() error {

	t.wsUrl = config.T.WsUrl
	t.okexSymbols = utils.ParseStringToArry(config.T.Symbol, ",")
	t.okexKlineTime = utils.ParseStringToArry(config.T.Kline, ",")
	t.okexDepth = utils.ParseStringToArry(config.T.Depth, ",")
	t.hbInterval = config.T.HeartBeat

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

		go t.readLoop(c, rgc, symbol)
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

func (t *OkExChange) readLoop(c *websocket.Conn, rgc chan int, symbol string) {

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

		//fmt.Println("ws receive message: ", string(message))

		// 去除心跳回应
		if len(message) == len(`{"event":"pong"}`) {
			continue
		}

		js, err := simplejson.NewJson(message)

		if err != nil {
			logs.Error("%s_%s sub ws parse json error:%s, json: %s", symbol, err.Error(), message)
			return
		}

		t.msgHandler(js)
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

func gzipDecode(in []byte) ([]byte, error) {

	reader := flate.NewReader(bytes.NewReader(in))

	defer reader.Close()

	return ioutil.ReadAll(reader)

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
