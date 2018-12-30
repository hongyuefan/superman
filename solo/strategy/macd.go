package strategy

import (
	"fmt"

	"sync"

	"github.com/hongyuefan/superman/protocol"
	"github.com/hongyuefan/superman/solo/database"
	"github.com/hongyuefan/superman/solo/skeleton"
)

var MAX_LEN int64 = 10

type MACD struct {
	EMA12 float64
	EMA26 float64
	DEA   float64
	DIF   float64
	Time  int64
}

type StratMacd struct {
	skl       *skeleton.Skeleton
	lock      sync.RWMutex
	chanClose chan bool
}

func NewStratMacd() *StratMacd {
	return &StratMacd{
		skl:       skeleton.NewSkeleton(),
		chanClose: make(chan bool, 0),
	}
}

func (s *StratMacd) Init() {

	s.skl.Init()

	return
}

func (s *StratMacd) SetMacd(kl protocol.KLineType, ema12, ema26, dea, dif, macd float64, time int64) error {
	switch kl {

	case protocol.SPIDER_TYPE_KLINE_15MIN:

		macd := &database.MACD_15Min{
			EMA12: ema12,
			EMA26: ema26,
			DEA:   dea,
			DIF:   dif,
			MACD:  macd,
			Time:  time,
		}
		if num, err := database.SetMACD_15Min(macd); err != nil || num == 0 {
			return fmt.Errorf("set macd error %v", kl)
		}
		return nil

	case protocol.SPIDER_TYPE_KLINE_HOUR:

		macd := &database.MACD_Hour{
			EMA12: ema12,
			EMA26: ema26,
			DEA:   dea,
			DIF:   dif,
			Time:  time,
		}
		if num, err := database.SetMACD_Hour(macd); err != nil || num == 0 {
			return fmt.Errorf("set macd error %v", kl)
		}
		return nil

	case protocol.SPIDER_TYPE_KLINE_DAY:

		macd := &database.MACD_Day{
			EMA12: ema12,
			EMA26: ema26,
			DEA:   dea,
			DIF:   dif,
			Time:  time,
		}
		if num, err := database.SetMACD_Day(macd); err != nil || num == 0 {
			return fmt.Errorf("set macd error %v", kl)
		}
		return nil
	}

	return fmt.Errorf("kline type not surpost %v", kl)
}
func (s *StratMacd) GetLastMacd(kl protocol.KLineType, offset int64) (EMA12, EMA26, DEA, DIF float64, Time int64, err error) {
	switch kl {

	case protocol.SPIDER_TYPE_KLINE_15MIN:

		macd := []database.MACD_15Min{}

		if num, err := database.GetMACD_15Min_Last(macd, 1, offset); err != nil || num == 0 {
			return 0, 0, 0, 0, 0, err
		}
		return macd[0].EMA12, macd[0].EMA26, macd[0].DEA, macd[0].DIF, macd[0].Time, nil

	case protocol.SPIDER_TYPE_KLINE_HOUR:

		macd := []database.MACD_Hour{}

		if num, err := database.GetMACD_Hour_Last(macd, 1, offset); err != nil || num == 0 {
			return 0, 0, 0, 0, 0, err
		}
		return macd[0].EMA12, macd[0].EMA26, macd[0].DEA, macd[0].DIF, macd[0].Time, nil

	case protocol.SPIDER_TYPE_KLINE_DAY:

		macd := []database.MACD_Day{}

		if num, err := database.GetMACD_Day_Last(macd, 1, offset); err != nil || num == 0 {
			return 0, 0, 0, 0, 0, err
		}
		return macd[0].EMA12, macd[0].EMA26, macd[0].DEA, macd[0].DIF, macd[0].Time, nil
	}

	return 0, 0, 0, 0, 0, fmt.Errorf("kline type not surpost %v", kl)
}

func (s *StratMacd) GetMacd(kl protocol.KLineType, time int64) (EMA12, EMA26, DEA, DIF float64, Time int64, err error) {

	switch kl {

	case protocol.SPIDER_TYPE_KLINE_15MIN:

		macd := &database.MACD_15Min{Time: time}

		if err := database.GetMACD_15Min(macd, "time"); err != nil {
			return 0, 0, 0, 0, 0, err
		}
		return macd.EMA12, macd.EMA26, macd.DEA, macd.DIF, macd.Time, nil

	case protocol.SPIDER_TYPE_KLINE_HOUR:

		macd := &database.MACD_Hour{Time: time}

		if err := database.GetMACD_Hour(macd, "time"); err != nil {
			return 0, 0, 0, 0, 0, err
		}
		return macd.EMA12, macd.EMA26, macd.DEA, macd.DIF, macd.Time, nil

	case protocol.SPIDER_TYPE_KLINE_DAY:

		macd := &database.MACD_Day{Time: time}

		if err := database.GetMACD_Day(macd, "time"); err != nil {
			return 0, 0, 0, 0, 0, err
		}
		return macd.EMA12, macd.EMA26, macd.DEA, macd.DIF, macd.Time, nil
	}

	return 0, 0, 0, 0, 0, fmt.Errorf("kline type not surpost %v", kl)

}

func (s *StratMacd) Calculation(kl protocol.KLineType) error {

	var (
		prema12 float64
		prema26 float64
		predea  float64
		err     error
	)
	//获取kline数据
	kls, ok := s.skl.GetKline(kl, 1)
	if !ok {
		return fmt.Errorf("kline %v no data", kl)
	}

	_, _, _, _, _, err = s.GetMacd(kl, kls[0].Time)
	if err != nil {
		prema12, prema26, predea, _, _, err = s.GetLastMacd(kl, 0)
		if err != nil {
			return err
		}
	}
	prema12, prema26, predea, _, _, err = s.GetLastMacd(kl, 1)
	if err != nil {
		return err
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

	fmt.Println("Calculation:", "EMA12:", curEMA12, "EMA26:", curEMA26, "DIF:", DIF, "DEA:", DEA, "TIME:", kls[0].Time)
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
	case protocol.NOTICE_KLINE_15MIN:
		//s.Calculation(protocol.SPIDER_TYPE_KLINE_15MIN)
		break
	case protocol.NOTICE_KLINE_HOUR:
		//s.Calculation(protocol.SPIDER_TYPE_KLINE_HOUR)
		break
	case protocol.NOTICE_KLINE_DAY:
		//s.Calculation(protocol.SPIDER_TYPE_KLINE_DAY)
		break
	}
}
