package main

import (
	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/utils"
)

func main() {

	utils.InitCnf()

	utils.InitLogger("solo", logs.LevelInfo)

	logs.Info("****************************************************")
	logs.Info("solo start...")
	logs.Info("  ")
	logs.Info("  ")
	logs.Info("  ")
	logs.Info("****************************************************")

	RunServer()
}
