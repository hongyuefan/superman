package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/hongyuefan/superman/config"
	"github.com/hongyuefan/superman/kfc"
	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/models"
	"github.com/hongyuefan/superman/protocol"
	"github.com/hongyuefan/superman/store"
	"github.com/hongyuefan/superman/utils"
)

func main() {

	utils.InitCnf()

	utils.InitLogger("store", logs.LevelInfo)

	logs.Info("****************************************************")
	logs.Info("storage start...")
	logs.Info("config file: ", config.T.CnfPath)

	if err := orm.RegisterDataBase("default", "mysql", config.T.SqlCon); err != nil {
		logs.Error("Connect SqlDB Error ", err.Error())
		return
	}

	conf, err := models.GetConfig()
	if err != nil {
		logs.Error("SqlDB GetConfig  Error ", err.Error())
		return
	}

	logs.Info("Store Path ", conf.Path)
	logs.Info("****************************************************")

	if err := StartServer(conf.Kafka, conf.Path, conf.ExChanges); err != nil {
		logs.Error("Start Server Error ", err.Error())
		os.Exit(-1)
	}
}

func StartServer(kafka, path, exchanges string) error {

	brokers := utils.ParseStringToArry(kafka, ",")

	topics := []string{protocol.TOPIC_OKEX_SPIDER_DATA}

	kfc.InitClient(brokers)

	err := kfc.TobeConsumer(topics)

	if err != nil {
		return fmt.Errorf("Init Kafka consumer error %v", err.Error())
	}

	logs.Info("connect to kafka broker [%s] ok ...", kafka)

	ch := make(chan int)

	server := store.NewStore(path, utils.ParseStringToArry(exchanges, ","), ch)

	if err := server.OnStart(); err != nil {
		return err
	}

	serverLoop(ch)

	return nil
}

func serverLoop(ch chan int) {

	defer close(ch)

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt)

	tc := time.NewTimer(time.Second)
	defer tc.Stop()

	for {
		select {
		case <-signals:
			logs.Info("recv a break signal, exit storage...")
			ch <- store.STG_CMD_EXIT
			<-time.After(time.Second)
			kfc.ExitConsumer()
			return

		case <-tc.C:
			tc.Reset(time.Second)
			if isEndOfDay() {
				ch <- store.STG_CMD_SWITCH_TRADINGDAY
			}
		}
	}
}

func isEndOfDay() bool {

	t := time.Now()

	if t.Hour() == 0 && t.Minute() == 0 && (t.Second() >= 0 || t.Second() < 2) {
		return true
	}

	return false
}
