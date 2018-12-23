package base

import (
	"github.com/hongyuefan/superman/config"
	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/utils"
	"github.com/okcoin-okex/open-api-v3-sdk/okex-go-sdk-api"
)

type BaseData struct {
	KLine   *KLineHandler
	Orders  *Order
	Wallets *Wallet
}

func NewBaseData() *BaseData {
	return &BaseData{
		KLine: NewKLineHandler(),
	}
}

func (b *BaseData) Init() {

	utils.InitCnf()

	utils.InitLogger("solo", logs.LevelInfo)

	client := okex.NewClient(okex.Config{ApiKey: config.T.ApiKey, SecretKey: config.T.ScretKey, Passphrase: config.T.PassPhrase, Endpoint: config.T.EndPoint, TimeoutSecond: 45, I18n: "en_US", IsPrint: false})

	b.Orders = NewOrder(client)
	b.Wallets = NewWallet(client)

	logs.Info("basedata init ok")
}
