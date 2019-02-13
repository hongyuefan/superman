package config

// app config
type AppCnf struct {
	LogPath string
	CnfPath string
	SqlCon  string

	//solo used
	ApiKey     string
	ScretKey   string
	PassPhrase string
	EndPoint   string
	Strategy   string
	Port       string

	AppID  string
	AppKey string
	TplId  int
	Mobile string
	Rates  string

	WsUrl     string
	Symbol    string
	Kline     string
	Depth     string
	HeartBeat int64
}

var T *AppCnf

func (c *AppCnf) LoadConfig(cnfPath string) (err error) {
	cnf, err := NewConfig("json", cnfPath)
	if err != nil {
		return err
	}
	c.CnfPath = cnfPath
	c.LogPath = cnf.String("log")
	c.SqlCon = cnf.String("sqlcon")
	c.Port = cnf.String("port")

	c.ApiKey = cnf.String("api_key")
	c.ScretKey = cnf.String("secret_key")
	c.PassPhrase = cnf.String("pass_phrase")
	c.EndPoint = cnf.String("end_point")
	c.Strategy = cnf.String("strategy")

	c.WsUrl = cnf.String("ws_url")
	c.Symbol = cnf.String("symbols")
	c.Kline = cnf.String("klines")
	c.Depth = cnf.String("depth")
	c.HeartBeat, _ = cnf.Int64("heart_beat")

	c.AppID = cnf.String("app_id")
	c.AppKey = cnf.String("app_key")
	c.TplId, _ = cnf.Int("tpl_id")
	c.Mobile = cnf.String("mobile")
	c.Rates = cnf.String("rates")

	return err
}

func newAppCnf() *AppCnf {
	return &AppCnf{}
}

func init() {
	T = newAppCnf()
}
