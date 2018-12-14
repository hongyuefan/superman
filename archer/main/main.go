package main

import (
	"fmt"
	"os"

	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/utils"
)

func main() {

	utils.InitCnf()

	utils.InitLogger("archer", logs.LevelInfo)

	if err := RunServer(); err != nil {
		fmt.Println(err)
		logs.Error("archer exit -1")
		os.Exit(-1)
	}
}
