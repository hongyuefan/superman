package store

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/hongyuefan/superman/kfc"
	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/models"
	"github.com/hongyuefan/superman/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

/*
 path的格式为：/usr/slash/data/
 data下是各个交易所名称，交易所下面是日期
 /usr/slash/data/okex/2017-12-04/quote
*/

const (
	STG_CMD_EXIT              = 1
	STG_CMD_SWITCH_TRADINGDAY = 2
)

type Store struct {
	DBName      string
	Path        string
	TradingDay  string
	ExChanges   []string
	Ch          chan int
	mDBs        map[string]*leveldb.DB
	mIndexCount map[string]uint64
}

func NewStore(path string, exchanges []string, ch chan int) *Store {
	return &Store{
		Path:        path,
		ExChanges:   exchanges,
		Ch:          ch,
		mDBs:        make(map[string]*leveldb.DB),
		mIndexCount: make(map[string]uint64),
		DBName:      "quotation",
	}
}

func (s *Store) loadIndex(exchange, tradingDay string) {

	v, err := models.GetStore(tradingDay, exchange)
	if err != nil {
		logs.Warn("loadIndex error :", err.Error())
	}
	if v != nil {
		s.mIndexCount[exchange] = v.ICount
	}
	s.mIndexCount[exchange] = 0
}

func (s *Store) saveIndex(exchange, tradingDay string, index uint64) {

	v, err := models.GetStore(tradingDay, exchange)
	if err != nil {
		logs.Warn("saveIndex error :", err.Error())
	}
	if v != nil {
		v.ICount = index
		if err := models.UpdateStore(v, "index_count"); err != nil {
			logs.Warn("saveIndex updateStore error :", err.Error())
		}
		return
	}
	if _, err := models.AddStore(&models.Store{ExChange: exchange, FileName: tradingDay, ICount: index}); err != nil {
		logs.Error("saveIndex addStore error :", err.Error())
	}
	return
}

func (s *Store) OnStart() error {

	s.TradingDay = getCurrDate()

	for _, exchange := range s.ExChanges {

		s.loadIndex(exchange, s.TradingDay)

		fileName := s.makeDBFileName(s.Path, exchange, s.TradingDay)

		db, err := leveldb.OpenFile(fileName, nil)

		if err != nil {
			logs.Error("open leveldb file error [%s]", err.Error())
			return err
		}

		s.mDBs[exchange] = db
	}

	go s.storeLoop()

	return nil
}

func getCurrDate() string {
	t := time.Now()
	return fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

func (s *Store) makeDBFileName(path string, exchange string, tradingDay string) string {
	return path + exchange + "/" + tradingDay + "/" + s.DBName
}

func (s *Store) storeLoop() {

	defer s.doStgExit()

	for {
		select {
		case msg := <-kfc.ReadMessages():

			s.handlerMsg(msg)

		case cmd, ok := <-s.Ch:
			if !ok || cmd == STG_CMD_EXIT {
				return
			}
			if cmd == STG_CMD_SWITCH_TRADINGDAY {
				if !s.SwitchTradingDay() {
					return
				}
			}
		}
	}
}

func (s *Store) doStgExit() {

	for exchange, v := range s.mDBs {

		s.saveIndex(exchange, s.TradingDay, s.mIndexCount[exchange])

		v.Close()
	}
}

func (s *Store) SwitchTradingDay() bool {

	newTradingDay := getCurrDate()

	if newTradingDay <= s.TradingDay {
		return true
	}

	oldTradingDay := s.TradingDay

	s.TradingDay = newTradingDay

	for exchange, v := range s.mDBs {

		v.Close()

		fileName := s.makeDBFileName(s.Path, exchange, s.TradingDay)

		db, err := leveldb.OpenFile(fileName, nil)

		if err != nil {
			logs.Error("store switch tradingday, openfile error: %s", err.Error())
			return false
		}

		s.mDBs[exchange] = db

		s.mIndexCount[exchange] = 0

		logs.Info("exchange [%s] has switch tradingDay [%s] -> [%s]", exchange, oldTradingDay, newTradingDay)
	}

	return true
}

func (s *Store) handlerMsg(msg *sarama.ConsumerMessage) bool {

	exchange := string(msg.Key)

	db, ok := s.mDBs[exchange]

	if !ok {
		logs.Error("store not support exchange: %s", exchange)
		return false
	}

	key := utils.UintTobytes(s.mIndexCount[exchange])

	err := db.Put(key, msg.Value, nil)

	if err != nil {
		logs.Error("store write [%s] error [%s]", s.makeDBFileName(s.Path, exchange, s.TradingDay), err.Error())
		return false
	}

	s.mIndexCount[exchange] += 1

	return true
}
