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

	c.ApiKey = cnf.String("api_key")
	c.ScretKey = cnf.String("secret_key")
	c.PassPhrase = cnf.String("pass_phrase")
	c.EndPoint = cnf.String("end_point")

	c.WsUrl = cnf.String("ws_url")
	c.Symbol = cnf.String("symbols")
	c.Kline = cnf.String("klines")
	c.Depth = cnf.String("depth")
	c.HeartBeat, _ = cnf.Int64("heart_beat")

	return err
}

func newAppCnf() *AppCnf {
	return &AppCnf{}
}

func init() {
	T = newAppCnf()
}
