package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/hongyuefan/superman/config"
	"github.com/hongyuefan/superman/krang"
	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/replay"
	"github.com/hongyuefan/superman/strategy/mavg"
	"github.com/hongyuefan/superman/utils"
)

func main() {

	utils.InitCnf()

	utils.InitLogger("replay", logs.LevelDebug)

	logs.Info("****************************************************")
	logs.Info("replay start...")
	logs.Info("  ")
	logs.Info("  ")
	logs.Info("  ")
	logs.Info("****************************************************")

	if err := orm.RegisterDataBase("default", "mysql", config.T.SqlCon); err != nil {
		logs.Error("Connect SqlDB Error ", err.Error())
		return
	}

	if err := RunServer(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func RunServer() error {

	mavg.RegisStrategy()

	ch := make(chan int)

	r, err := replay.StartReplay(ch)

	if err != nil {
		return err
	}

	krang.SetKrangReplay(r)

	err = krang.StartKrang(ch, true)
	if err != nil {
		return err
	}

	serverLoop(ch)
	return nil
}

func serverLoop(ch chan int) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case <-signals:
			logs.Info("recv a break signal, exit replay ...")
			<-time.After(3 * time.Second)
			return

		case <-ch:
			logs.Info("replay is all done. ")
			return
		}
	}
}
