package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/hongyuefan/superman/config"
	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/solo/database"
	"github.com/hongyuefan/superman/solo/strategy"
)

type StrateGY interface {
	Init()
	OnTicker()
	OnClose()
}

func RunServer() {

	database.RegistDB()

	var strat StrateGY

	switch config.T.Strategy {
	case "macd":
		strat = strategy.NewStratKDJ()
		break
	case "kdj":
		strat = strategy.NewStratMacd()
		break
	}

	strat.Init()

	strat.OnTicker()

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt)

	<-signals

	logs.Info("recv a break signal, exit solo...")

	strat.OnClose()

	<-time.After(time.Second)

}
