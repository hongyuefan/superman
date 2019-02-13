package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	gin "github.com/gin-gonic/gin"
	"github.com/hongyuefan/superman/config"
	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/solo/database"
	"github.com/hongyuefan/superman/solo/strategy"
)

type StrateGY interface {
	Init()
	OnTicker()
	OnClose()
	Handler(*gin.Context)
}

func ServerHandler(s StrateGY) {

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.GET("/handler", s.Handler)
	}

	fmt.Println("Listen:", config.T.Port)

	http.ListenAndServe(":"+config.T.Port, router)
}

func RunServer() {

	database.RegistDB()

	var strat StrateGY

	switch config.T.Strategy {
	case "updown":
		strat = strategy.NewStratUpDown()
	case "kdj":
		strat = strategy.NewStratKDJ()
		break
	case "macd":
		strat = strategy.NewStratMacd()
		break
	case "sample":
		strat = strategy.NewSampleSt()
		break
	}

	strat.Init()

	go ServerHandler(strat)

	strat.OnTicker()

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt)

	<-signals

	logs.Info("recv a break signal, exit solo...")

	strat.OnClose()

	<-time.After(time.Second)

}
