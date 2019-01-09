package strategy

import (
	"fmt"
	"sync"

	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/protocol"
	"github.com/hongyuefan/superman/solo/database"
	"github.com/hongyuefan/superman/solo/skeleton"
)

type KDJ struct {
	K   float64
	D   float64
	J   float64
	RSV float64
}

type StratKDJ struct {
	skl       *skeleton.Skeleton
	lock      sync.RWMutex
	chanClose chan bool
	mapFlag   map[protocol.KLineType]bool
}

func NewStratKDJ() *StratKDJ {
	return &StratKDJ{
		skl:       skeleton.NewSkeleton(),
		chanClose: make(chan bool, 0),
		mapFlag: map[protocol.KLineType]bool{
			protocol.SPIDER_TYPE_KLINE_5MIN:  false,
			protocol.SPIDER_TYPE_KLINE_15MIN: false,
			protocol.SPIDER_TYPE_KLINE_HOUR:  false,
			protocol.SPIDER_TYPE_KLINE_DAY:   false,
		},
	}
}

func (s *StratKDJ) OnInit() {

	s.skl.Init()

	return
}

func (s *StratKDJ) OnTicker() {

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

func (s *StratKDJ) OnClose() {

	s.skl.Close()

	close(s.chanClose)

}

func (s *StratKDJ) SetKDJ(kl protocol.KLineType, K, D, J, RSV, time int64) error {

	switch kl {

	case protocol.SPIDER_TYPE_KLINE_5MIN:

		kdj := &database.KDJ_5Min{
			K:    K,
			D:    D,
			J:    J,
			RSV:  RSV,
			Time: time,
		}
		if num, err := database.SetKDJ_5Min(kdj); err != nil || num == 0 {
			return fmt.Errorf("set kdj error %v", kl)
		}
		return nil
	}
	return fmt.Errorf("kline type not surpost %v", kl)
}

func (s *StratKDJ) GetLastKDJ(kl protocol.KLineType, offset int64) (k, d, j, rsv float64, Time int64, err error) {
	switch kl {

	case protocol.SPIDER_TYPE_KLINE_5MIN:

		kdj := []database.KDJ_5Min{}

		if _, err := database.GetKDJ_5Min_Last(&kdj, 1, offset); err != nil || len(kdj) == 0 {
			return 0, 0, 0, 0, 0, fmt.Errorf("GetLastKDJ Error ")
		}
		return kdj[0].K, kdj[0].D, kdj[0].J, kdj[0].RSV, macd[0].Time, nil

	}
	return 0, 0, 0, 0, 0, fmt.Errorf("kline type not surpost %v", kl)
}

func (s *StratKDJ) GetKDJ(kl protocol.KLineType, time int64) (k, d, j, rsv float64, Time int64, err error) {

	switch kl {

	case protocol.SPIDER_TYPE_KLINE_5MIN:

		kdj := []database.KDJ_5Min{}

		if err := database.GetKDJ_5Min(kdj, "time"); err != nil {
			return 0, 0, 0, 0, err
		}
		return kdj.K, kdj.D, kdj.j, kdj.RSV, kdj.Time, nil
	}
	return 0, 0, 0, 0, 0, fmt.Errorf("kline type not surpost %v", kl)

}

func (s *StratKDJ) Calculation(kl protocol.KLineType) error {

	var (
		preK    float64
		preD    float64
		preJ    float64
		preRSV  float64
		pretime int64
		err     error
	)
	//获取kline数据
	kls, ok := s.skl.GetKline(kl, 1)
	if !ok {
		return fmt.Errorf("kline %v no data", kl)
	}

	_, _, _, _, _, err = s.GetKDJ(kl, kls[0].Time)
	if err != nil {
		preK, preD, preJ, _, pretime, err = s.GetLastKDJ(kl, 0)
		if err != nil {
			return err
		}
		s.doOrder(kl, premacd, pretime)
	} else {
		preK, preD, preJ, _, _, err = s.GetLastMacd(kl, 1)
		if err != nil {
			return err
		}
	}

	//计算ema指数
	curEMA12 := prema12*11/13 + kls[0].Close*2/13
	curEMA26 := prema26*25/27 + kls[0].Close*2/27

	//计算当前 dif、dea 指标
	DIF := curEMA12 - curEMA26
	DEA := predea*8/10 + DIF*2/10
	MACD := (DIF - DEA) * 2

	if err := s.SetMacd(kl, curEMA12, curEMA26, DEA, DIF, MACD, kls[0].Time); err != nil {
		return err
	}

	fmt.Println("Calculation:", kl, "EMA12:", curEMA12, "EMA26:", curEMA26, "DIF:", DIF, "DEA:", DEA, "MACD:", MACD, "TIME:", kls[0].Time)

	return nil
}

func (s *StratKDJ) dispatchMsg(symbol string, notice protocol.NoticeType) {

	switch notice {

	case protocol.NOTICE_KLINE_5MIN:

		if !s.mapFlag[protocol.SPIDER_TYPE_KLINE_5MIN] {

			s.judgeKDJ(protocol.SPIDER_TYPE_KLINE_5MIN)

		}

		s.Calculation(protocol.SPIDER_TYPE_KLINE_5MIN)

		break

	case protocol.NOTICE_KLINE_15MIN:
		break
	case protocol.NOTICE_KLINE_HOUR:
		break
	case protocol.NOTICE_KLINE_DAY:
		break
	}
}

func (s *StratKDJ) judgeKDJ(kl protocol.KLineType) {

}
