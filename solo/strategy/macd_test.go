package strategy

import (
	"testing"

	"github.com/hongyuefan/superman/protocol"
)

func TestBench(t *testing.T) {

	std := NewStratMacd()

	std.GetLastMacd(protocol.SPIDER_TYPE_KLINE_5MIN, 0)
}
