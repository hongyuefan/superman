package strategy

import (
	"fmt"
	"sync"

	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/protocol"
	"github.com/hongyuefan/superman/solo/base"
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

func (s *StratKDJ) Init() {

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

func (s *StratKDJ) SetKDJ(kl protocol.KLineType, K, D, J, RSV float64, time int64) error {

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

		var kdj []database.KDJ_5Min

		if _, err := database.GetKDJ_5Min_Last(&kdj, 1, offset); err != nil || len(kdj) == 0 {
			return 0, 0, 0, 0, 0, fmt.Errorf("GetLastKDJ Error %v %v", kl, offset)
		}
		return kdj[0].K, kdj[0].D, kdj[0].J, kdj[0].RSV, kdj[0].Time, nil

	}
	return 0, 0, 0, 0, 0, fmt.Errorf("kline type not surpost %v", kl)
}

func (s *StratKDJ) GetKDJ(kl protocol.KLineType, time int64) (k, d, j, rsv float64, Time int64, err error) {

	switch kl {

	case protocol.SPIDER_TYPE_KLINE_5MIN:

		var kdj database.KDJ_5Min

		if err := database.GetKDJ_5Min(&kdj, "time"); err != nil {
			return 0, 0, 0, 0, 0, err
		}
		return kdj.K, kdj.D, kdj.J, kdj.RSV, kdj.Time, nil
	}
	return 0, 0, 0, 0, 0, fmt.Errorf("kline type not surpost %v", kl)

}

func (s *StratKDJ) doOrder(kl protocol.KLineType, preK, curK float64) {

	if (curK <= 20 && curK > preK) || (curK > 50 && curK <= 80 && curK > preK) {
		fmt.Println("======Buy=====:", kl, s.skl.GetTicker().Last, curK-preK)
	}

	if (curK >= 0 && curK < preK) || (curK > 20 && curK <= 50 && curK < preK) {
		fmt.Println("======Sell=====:", kl, s.skl.GetTicker().Last, curK-preK)
	}
	return
}

func (s *StratKDJ) Calculation(kl protocol.KLineType) error {

	var (
		preK   float64
		preD   float64
		curRSV float64
		err    error
	)

	//获取kline数据
	kls, ok := s.skl.GetKline(kl, 10)

	if !ok || len(kls) < 10 {
		return fmt.Errorf("kline %v not enough data", kl)
	}

	k, d, j, rsv, t, err := s.GetKDJ(kl, kls[0].Time)

	fmt.Println("getkdj:", k, d, j, rsv, t, kls[0].Time, err)

	if err != nil {

		preK, preD, _, _, _, err = s.GetLastKDJ(kl, 0)

		if err != nil {
			return err
		}

		curRSV = s.rsv(kls[0].Close, kls[1:])

		curK := 2/3*preK + 1/3*curRSV

		curD := 2/3*preD + 1/3*curK

		curJ := 3*curK - 2*curD

		if err := s.SetKDJ(kl, curK, curD, curJ, curRSV, kls[0].Time); err != nil {
			return err
		}

		fmt.Println("Calculation:", kl, "K:", curK, "D:", curD, "J:", curJ, "RSV:", curRSV, "TIME:", kls[0].Time)

		s.doOrder(kl, preK, curK)

	} else {

		preK, preD, _, _, _, err = s.GetLastKDJ(kl, 1)

		if err != nil {
			return err
		}

		curRSV = s.rsv(kls[0].Close, kls[:9])

		curK := 2/3*preK + 1/3*curRSV

		curD := 2/3*preD + 1/3*curK

		curJ := 3*curK - 2*curD

		if err := s.SetKDJ(kl, curK, curD, curJ, curRSV, kls[0].Time); err != nil {
			return err
		}

		fmt.Println("Calculation:", kl, "K:", curK, "D:", curD, "J:", curJ, "RSV:", curRSV, "TIME:", kls[0].Time)
	}

	return nil
}

func (s *StratKDJ) rsv(c float64, arrys []base.KLineDetail) float64 {

	var (
		low  float64 = 999999
		high float64 = 0
	)

	for _, d := range arrys {
		if d.Low < low {
			low = d.Low
		}
		if d.High > high {
			high = d.High
		}
	}

	fmt.Println("high:", high, "low:", low)

	return (c - low) / (high - low) * 100
}

func (s *StratKDJ) dispatchMsg(symbol string, notice protocol.NoticeType) {

	switch notice {

	case protocol.NOTICE_KLINE_5MIN:

		if !s.mapFlag[protocol.SPIDER_TYPE_KLINE_5MIN] {
			s.judgeKDJ(protocol.SPIDER_TYPE_KLINE_5MIN)
		}

		if err := s.Calculation(protocol.SPIDER_TYPE_KLINE_5MIN); err != nil {
			fmt.Println("calculate error:", err)
		}

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

	var err error

	_, ok := s.skl.GetKline(kl, 1)

	if !ok {
		logs.Error("kline %v no data", kl)
		return
	}

	_, _, _, _, _, err = s.GetLastKDJ(kl, 0)

	if err != nil {

		if err = s.SetKDJ(kl, 50, 50, 50, 0, 0); err != nil {
			logs.Error("kline %v setKDJ error:%v", kl, err)
			return
		}
	}

	s.mapFlag[kl] = true

	return
}
