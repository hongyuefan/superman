package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/hongyuefan/superman/config"
	"github.com/hongyuefan/superman/kfc"
	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/models"
	exs "github.com/hongyuefan/superman/spider/exchanges"
	"github.com/hongyuefan/superman/utils"
)

func RunServer() error {

	err := orm.RegisterDataBase("default", "mysql", config.T.SqlCon)
	if err != nil {
		return err
	}

	conf, err := models.GetConfig()
	if err != nil {
		logs.Error("Connect SqlDB Error : %v", err.Error())
		return err
	}

	brokers := utils.ParseStringToArry(conf.Kafka, ",")

	exchanges := utils.ParseStringToArry(conf.ExChanges, ",")

	kfc.InitClient(brokers)

	err = kfc.TobeProducer()
	if err != nil {
		logs.Error("InitKafkaClient producer error ", err.Error())
		return err
	}

	logs.Info("connect to kafka broker ok ...")

	if err := exs.StartExchange(exchanges); err != nil {
		return err
	}

	serverLoop()
	return nil
}

func serverLoop() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	<-signals
	logs.Info("recv a break signal, exit spider...")
	kfc.ExitProducer()

	<-time.After(time.Second)
}
