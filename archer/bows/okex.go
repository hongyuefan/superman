package bows

import (
	"encoding/json"

	"github.com/hongyuefan/superman/logs"
	"github.com/hongyuefan/superman/models"
	"github.com/hongyuefan/superman/protocol"
	"github.com/hongyuefan/superman/utils"
	"github.com/okcoin-okex/open-api-v3-sdk/okex-go-sdk-api"
)

type okexArcher struct {
	restUrl string
	errm    map[int]string
}

func newOkexArcher() Archer {
	return &okexArcher{
		errm: make(map[int]string),
	}
}

func (t *okexArcher) Init() error {

	spider, err := models.GetSpiderByName("okex")

	if err != nil {
		return err
	}

	t.restUrl = spider.RpcUrl

	utils.InitOkexErrorMap(t.errm)

	return nil
}

func (t *okexArcher) GetOkClient(cmd *ArcherCmd) (*okex.Client, error) {

	confParam := &protocol.OkexUserParam{}

	if err := json.Unmarshal(cmd.PayLoad, confParam); err != nil {
		return nil, err
	}

	return t.NewOkClient(confParam.ApiKey, confParam.SecretKey, confParam.PassPhrase), nil
}

func (t *okexArcher) NewOkClient(apiKey, secretKey, passPhrase string) *okex.Client {
	return okex.NewClient(okex.Config{ApiKey: apiKey, SecretKey: secretKey, Passphrase: passPhrase, Endpoint: t.restUrl, TimeoutSecond: 45, I18n: "en_US", IsPrint: false})
}

func (t *okexArcher) Run(archerCh chan *ArcherCmd) {

	for {
		cmd := <-archerCh

		logs.Info("okex receive data:", cmd.ToString())

		switch cmd.BusiId {

		case protocol.CMD_EXIST:
			logs.Info("okex archer exit")
			return

		case protocol.CMD_QRY_ACCOUNTS:
			t.queryAccounts(cmd)

		case protocol.CMD_QRY_ACCOUNT:
			t.queryAccount(cmd)

		case protocol.CMD_DO_ORDER:
			t.doOrder(cmd)

		case protocol.CMD_CANCEL_ORDER:
			t.canselOrder(cmd)

		case protocol.CMD_QRY_PENDING_ORDER:
			t.queryPendingOrders(cmd)

		case protocol.CMD_QRY_ORDER:
			t.queryOrder(cmd)

		}

		logs.Info("okex archer excute [%d] success", cmd.BusiId)
	}
}

////////////////////////////////////////////////////////////////////////////////

func okexArcherReply(busiId int32, reqSerial uint32, msg interface{}) error {
	byt, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return utils.PackAndReplyToBroker(protocol.TOPIC_OKEX_ARCHER_RSP, "okex", int(busiId), string(byt))
}

////////////////////////////////////////////////////////////////////////////////

func (t *okexArcher) queryAccounts(cmd *ArcherCmd) {

	client, err := t.GetOkClient(cmd)
	if err != nil {
		logs.Warn("queryAccounts GetOkClient error:%s", err.Error())
		return
	}
	result, err := client.SpotGetAccounts()
	if err != nil {
		logs.Warn("queryAccounts SpotGetAccounts error:%s", err.Error())
		return
	}
	if err := okexArcherReply(cmd.BusiId, cmd.ReqSerial, result); err != nil {
		logs.Warn("queryAccounts reply error:%s", err.Error())
	}
	return
}

func (t *okexArcher) queryAccount(cmd *ArcherCmd) {

	client, err := t.GetOkClient(cmd)
	if err != nil {
		logs.Warn("queryAccount GetOkClient error:%s", err.Error())
		return
	}

	param := new(okex.SpotAccountParams)

	if err := json.Unmarshal(cmd.PayLoad, param); err != nil {
		logs.Warn("queryAccount unmarshal error:%s", err.Error())
		return
	}

	result, err := client.SpotGetAccountCurrency(*param)
	if err != nil {
		logs.Warn("queryAccount SpotGetAccountCurrency error:%s", err.Error())
		return
	}
	if err := okexArcherReply(cmd.BusiId, cmd.ReqSerial, result); err != nil {
		logs.Warn("queryAccount reply error:%s", err.Error())
	}
	return
}

func (t *okexArcher) doOrder(cmd *ArcherCmd) {

	client, err := t.GetOkClient(cmd)
	if err != nil {
		logs.Warn("doOrder GetOkClient error:%s", err.Error())
		return
	}

	param := new(okex.SpotOrderParams)

	if err := json.Unmarshal(cmd.PayLoad, param); err != nil {
		logs.Warn("doOrder unmarshal error:%s", err.Error())
		return
	}

	result, err := client.SpotDoOrder(*param)
	if err != nil {
		logs.Warn("doOrder SpotDoOrder error:%s", err.Error())
		return
	}
	if err := okexArcherReply(cmd.BusiId, cmd.ReqSerial, result); err != nil {
		logs.Warn("doOrder reply error:%s", err.Error())
	}
	return
}

func (t *okexArcher) canselOrder(cmd *ArcherCmd) {

	client, err := t.GetOkClient(cmd)
	if err != nil {
		logs.Warn("canselOrder GetOkClient error:%s", err.Error())
		return
	}

	param := new(okex.SpotCanselOrderParams)

	if err := json.Unmarshal(cmd.PayLoad, param); err != nil {
		logs.Warn("canselOrder Unmarshal error:%s", err.Error())
		return
	}

	result, err := client.SpotCanselOrder(*param)
	if err != nil {
		logs.Warn("canselOrder SpotCanselOrder error:%s", err.Error())
		return
	}
	if err := okexArcherReply(cmd.BusiId, cmd.ReqSerial, result); err != nil {
		logs.Warn("canselOrder reply error:%s", err.Error())
	}
	return
}

func (t *okexArcher) queryPendingOrders(cmd *ArcherCmd) {

	client, err := t.GetOkClient(cmd)
	if err != nil {
		logs.Warn("queryPendingOrders GetOkClient error:%s", err.Error())
		return
	}

	param := new(okex.SpotGetPendingOrderParams)

	if err := json.Unmarshal(cmd.PayLoad, param); err != nil {
		logs.Warn("queryPendingOrders Unmarshal error:%s", err.Error())
		return
	}

	result, err := client.SpotGetPendingOrders(*param)
	if err != nil {
		logs.Warn("queryPendingOrders SpotGetPendingOrders ")
		return
	}
	if err := okexArcherReply(cmd.BusiId, cmd.ReqSerial, result); err != nil {
		logs.Warn("queryPendingOrders reply error:%s", err.Error())
	}
	return
}

func (t *okexArcher) queryOrder(cmd *ArcherCmd) {

	client, err := t.GetOkClient(cmd)
	if err != nil {
		logs.Warn("queryOrder GetOkClient error:%s", err.Error())
		return
	}
	param := new(okex.SpotGetOrderParams)

	if err := json.Unmarshal(cmd.PayLoad, param); err != nil {
		logs.Warn("queryOrder unmarshal error:%s", err.Error())
		return
	}
	result, err := client.SpotGetOrder(*param)
	if err != nil {
		logs.Warn("queryOrder SpotGetOrder error:%s", err.Error())
		return
	}
	if err := okexArcherReply(cmd.BusiId, cmd.ReqSerial, result); err != nil {
		logs.Warn("queryOrder reply error:%s", err.Error())
	}
	return
}
