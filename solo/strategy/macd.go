package strategy

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/hongyuefan/superman/config"
	"github.com/hongyuefan/superman/protocol"
	"github.com/hongyuefan/superman/solo/skeleton"
	"github.com/hongyuefan/superman/utils"
)

const (
	MAX_LENGTH = 10
)

var mapIntervel = map[protocol.KLineType]int64{
	protocol.SPIDER_TYPE_KLINE_1MIN:  60000,
	protocol.SPIDER_TYPE_KLINE_5MIN:  300000,
	protocol.SPIDER_TYPE_KLINE_15MIN: 900000,
	protocol.SPIDER_TYPE_KLINE_30MIN: 1800000,
	protocol.SPIDER_TYPE_KLINE_HOUR:  3600000,
	protocol.SPIDER_TYPE_KLINE_DAY:   86400000,
}

var mapPreEMA = map[protocol.KLineType]EMA{
	protocol.SPIDER_TYPE_KLINE_1MIN:  EMA{EMA12: 10, EMA26: 12},
	protocol.SPIDER_TYPE_KLINE_5MIN:  EMA{EMA12: 10, EMA26: 12},
	protocol.SPIDER_TYPE_KLINE_15MIN: EMA{EMA12: 10, EMA26: 12},
	protocol.SPIDER_TYPE_KLINE_30MIN: EMA{EMA12: 10, EMA26: 12},
	protocol.SPIDER_TYPE_KLINE_HOUR:  EMA{EMA12: 10, EMA26: 12},
	protocol.SPIDER_TYPE_KLINE_DAY:   EMA{EMA12: 10, EMA26: 12},
}

var mapPreMACD = map[protocol.KLineType]MACD{
	protocol.SPIDER_TYPE_KLINE_1MIN:  MACD{DEA: 10, DIF: 12, stamp: 10},
	protocol.SPIDER_TYPE_KLINE_5MIN:  MACD{DEA: 10, DIF: 12, stamp: 10},
	protocol.SPIDER_TYPE_KLINE_15MIN: MACD{DEA: 10, DIF: 12, stamp: 10},
	protocol.SPIDER_TYPE_KLINE_30MIN: MACD{DEA: 10, DIF: 12, stamp: 10},
	protocol.SPIDER_TYPE_KLINE_HOUR:  MACD{DEA: 10, DIF: 12, stamp: 10},
	protocol.SPIDER_TYPE_KLINE_DAY:   MACD{DEA: 10, DIF: 12, stamp: 10},
}

type MACD struct {
	DEA   float32
	DIF   float32
	stamp int64
}

type EMA struct {
	EMA12 float32
	EMA26 float32
}

type StratMacd struct {
	skl                *skeleton.Skeleton
	map_symbol_history map[string]MapHistory
	lock               sync.RWMutex
	chanClose          chan bool
}

func NewStratMacd() *StratMacd {
	return &StratMacd{
		skl:                skeleton.NewSkeleton(),
		chanClose:          make(chan bool, 0),
		map_symbol_history: make(map[string]MapHistory),
	}
}

type MapHistory map[protocol.KLineType]*History

type History struct {
	preEMA   EMA
	nowEMA   EMA
	map_MACD map[int]*MACD //历史macd数据值 ，1：macd
	lock     sync.RWMutex
}

func NewHistory() *History {
	return &History{
		map_MACD: make(map[int]*MACD),
	}
}

//清洗MACD数据
func (h *History) CleanMACD(intervel int64, curEMA EMA, macd *MACD) {
	h.lock.Lock()
	defer h.lock.Unlock()

	curMacd, ok := h.map_MACD[0]

	if !ok {
		h.map_MACD[0] = macd
		return
	}

	if (macd.stamp - curMacd.stamp) >= intervel {

		length := len(h.map_MACD)

		if length < MAX_LENGTH {
			for index := length - 1; index >= 0; index-- {
				h.map_MACD[index+1] = h.map_MACD[index]
			}
		} else {
			for index := length - 1; index >= 1; index-- {
				h.map_MACD[index] = h.map_MACD[index-1]
			}
		}

		h.preEMA = h.nowEMA
	}
	h.map_MACD[0] = macd
	h.nowEMA = curEMA
}

func (h *History) String() string {

	var macds []MACD

	for _, macd := range h.map_MACD {
		macds = append(macds, *macd)
	}
	return fmt.Sprintf("preEMA:%v,nowEMA:%v,MACD_length:%v,MACD:%v", h.preEMA, h.nowEMA, len(h.map_MACD), macds)
}

func (s *StratMacd) Init() {

	s.skl.Init()

	symbols := utils.ParseStringToArry(config.T.Symbol, ",")

	kls := []protocol.KLineType{protocol.SPIDER_TYPE_KLINE_1MIN, protocol.SPIDER_TYPE_KLINE_5MIN, protocol.SPIDER_TYPE_KLINE_15MIN, protocol.SPIDER_TYPE_KLINE_30MIN, protocol.SPIDER_TYPE_KLINE_HOUR, protocol.SPIDER_TYPE_KLINE_DAY}

	for _, symbol := range symbols {

		for _, kl := range kls {

			s.map_symbol_history[symbol] = make(map[protocol.KLineType]*History)

			m_kl_his := s.map_symbol_history[symbol]

			m_kl_his[kl] = NewHistory()

			his := m_kl_his[kl]

			his.preEMA = mapPreEMA[kl]

			his.map_MACD[0] = &MACD{DEA: mapPreMACD[kl].DEA}
		}
	}
}

func (s *StratMacd) Calculation(symbol string, kl protocol.KLineType) error {

	//获取symbol下的数据集合
	m_kl_his, ok := s.map_symbol_history[symbol]
	if !ok {
		return fmt.Errorf("symbol %s not support", symbol)
	}

	//获取K线指标下的相应EMA数据结构
	his, ok := m_kl_his[kl]
	if !ok {
		return fmt.Errorf("kline %v not support", kl)
	}

	//获取ticker数据
	ticker, ok := s.skl.GetTicker(symbol)
	if !ok {
		return fmt.Errorf("ticker %s no data", symbol)
	}

	fLast, err := strconv.ParseFloat(ticker.Last, 32)
	if err != nil {
		return fmt.Errorf("parseFloat ticker.last error :%s", err.Error())
	}

	//计算ema指数
	curEMA12 := his.preEMA.EMA12*11/13 + float32(fLast)*2/13
	curEMA26 := his.preEMA.EMA26*25/27 + float32(fLast)*2/27

	//计算当前 dif、dea 指标
	DIF := curEMA12 - curEMA26
	DEA := his.map_MACD[0].DEA*8/10 + DIF*2/10

	curMacd := &MACD{DIF: DIF, DEA: DEA, stamp: ticker.TimeStamp}

	his.CleanMACD(mapIntervel[kl], EMA{EMA12: curEMA12, EMA26: curEMA26}, curMacd)

	fmt.Println("Calculation Symbol", symbol, "kline type", kl.String(), "content:", his.String())

	return nil
}

func (s *StratMacd) OnTicker() {

	go s.skl.Run()

	for {
		select {
		case notice := <-s.skl.ChanNotice():
			s.dispatchMsg(notice.Symbol, notice.Notice)
		case <-s.chanClose:
			return
		}
	}
}

func (s *StratMacd) OnClose() {
	s.skl.Close()
	close(s.chanClose)
}

func (s *StratMacd) dispatchMsg(symbol string, notice protocol.NoticeType) {

	switch notice {
	case protocol.NOTICE_KLINE_1MIN:
		s.Calculation(symbol, protocol.SPIDER_TYPE_KLINE_1MIN)
		break
	case protocol.NOTICE_KLINE_5MIN:
		s.Calculation(symbol, protocol.SPIDER_TYPE_KLINE_5MIN)
		break
	case protocol.NOTICE_KLINE_15MIN:
		s.Calculation(symbol, protocol.SPIDER_TYPE_KLINE_15MIN)
		break
	case protocol.NOTICE_KLINE_30MIN:
		s.Calculation(symbol, protocol.SPIDER_TYPE_KLINE_30MIN)
		break
	case protocol.NOTICE_KLINE_HOUR:
		s.Calculation(symbol, protocol.SPIDER_TYPE_KLINE_HOUR)
		break
	case protocol.NOTICE_KLINE_DAY:
		s.Calculation(symbol, protocol.SPIDER_TYPE_KLINE_DAY)
		break
	case protocol.NOTICE_KLINE_WEEK:
		s.Calculation(symbol, protocol.SPIDER_TYPE_KLINE_WEEK)
		break
	}
}
