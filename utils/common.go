package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hongyuefan/superman/kfc"
	"github.com/hongyuefan/superman/message"
	"github.com/hongyuefan/superman/protocol"
)

func PackAndReplyToBroker(topic, key string, typ int, data string) error {

	msg := message.Messages{
		Type:  typ,
		Datas: data,
	}

	byt, _ := json.Marshal(msg)

	kfc.SendMessage(topic, key, byt)

	return nil
}

func MakeupSinfo(ex string, symbol string, contractType string) string {
	return ex + "_" + symbol + "_" + contractType
}

func UintTobytes(i uint64) []byte {
	return []byte(fmt.Sprintf("%d", i))
}

func BytesToUint(b []byte) uint64 {
	u, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		u = 0
	}
	return u
}

func KLineStr(kl int32) string {
	switch kl {
	case protocol.KL1Min:
		return "KL1Min"
	case protocol.KL3Min:
		return "KL3Min"
	case protocol.KL5Min:
		return "KL5Min"
	case protocol.KL15Min:
		return "KL15Min"
	case protocol.KL30Min:
		return "KL30Min"
	case protocol.KL1H:
		return "KL1H"
	case protocol.KL1D:
		return "KL1D"
	}
	return "未知"
}

func TSStr(ts int64) string {
	return time.Unix(ts, 0).Format(protocol.TM_LAYOUT_STR)
}

func IsZero32(f float32) bool {
	return f >= -0.000001 && f <= 0.000001
}

func IsZero64(f float64) bool {
	return f >= -0.000001 && f <= 0.000001
}
