package main

import (
	"fmt"
	"os"

	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/utils"
)

func main() {
	utils.InitCnf()
	utils.InitLogger("spider", logs.LevelInfo)

	logs.Info("****************************************************")
	logs.Info("spider start...")
	logs.Info("  ")
	logs.Info("  ")
	logs.Info("  ")
	logs.Info("****************************************************")

	if err := RunServer(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
