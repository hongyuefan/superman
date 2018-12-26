package base

import (
	"testing"

	"github.com/hongyuefan/superman/protocol"
)

func TestBench(t *testing.T) {

	symbol := "btc-usdt"

	handle := NewKLineHandler()

	handle.Handler(protocol.SPIDER_TYPE_KLINE_1MIN, symbol, "0", "1.1", "1.9", "1.0", "1.5", "1.4")

	handle.Handler(protocol.SPIDER_TYPE_KLINE_1MIN, symbol, "1", "1.1", "1.9", "1.0", "1.5", "1.4")

	handle.Handler(protocol.SPIDER_TYPE_KLINE_1MIN, symbol, "62", "1.1", "1.9", "1.0", "1.5", "1.4")

	handle.Handler(protocol.SPIDER_TYPE_KLINE_1MIN, symbol, "124", "1.1", "1.9", "1.0", "1.5", "1.4")

	handle.Handler(protocol.SPIDER_TYPE_KLINE_1MIN, symbol, "1240", "1.1", "1.9", "1.0", "1.5", "1.4")

	detal, ok := handle.Get(symbol, protocol.SPIDER_TYPE_KLINE_1MIN, 0)

	t.Log(detal, ok)

}
