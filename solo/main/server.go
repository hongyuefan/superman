package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/solo/database"
	"github.com/hongyuefan/superman/solo/strategy"
)

func RunServer() {

	database.RegistDB()

	strat := strategy.NewStratMacd()

	strat.Init()

	strat.OnTicker()

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt)

	<-signals

	logs.Info("recv a break signal, exit solo...")

	strat.OnClose()

	<-time.After(time.Second)

}
