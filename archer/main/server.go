package main

import (
	"github.com/astaxie/beego/orm"
	"github.com/hongyuefan/superman/archer/bows"
	"github.com/hongyuefan/superman/config"
	"github.com/hongyuefan/superman/models"
	"github.com/hongyuefan/superman/utils"
)

func RunServer() error {

	if err := orm.RegisterDataBase("default", "mysql", config.T.SqlCon); err != nil {
		return err
	}

	conf, err := models.GetConfig()
	if err != nil {
		return err
	}

	if err := bows.InitKafkaClient(utils.ParseStringToArry(conf.Kafka, ",")); err != nil {
		return err
	}

	bl := bows.InitBows()

	if err := bows.StartExArcher(utils.ParseStringToArry(conf.ExChanges, ","), bl); err != nil {
		return err
	}

	bows.StartCmdLoop(bl)

	return nil
}
