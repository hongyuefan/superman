package exchanges

import "github.com/hongyuefan/superman/logs"

/*
  交易所行情服务接口
*/
type Exchange interface {
	Init() error
	Run()
}

/*
  StartQuoters -- 启动行情拉取服务
*/
func StartExchange(exchanges []string) error {

	for _, ex := range exchanges {

		q := createExchange(ex)

		if q == nil {
			logs.Error("exchange [%s] is not exist")
			continue
		}

		if err := q.Init(); err != nil {
			logs.Error("exchange [%s] init fail, error:", ex, err.Error())
			return err
		}

		go q.Run()
	}

	return nil
}

func createExchange(exName string) Exchange {

	switch exName {
	case "okex":
		return NewOkExChange()
	}

	return nil
}
