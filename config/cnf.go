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

	return err
}

func newAppCnf() *AppCnf {
	return &AppCnf{}
}

func init() {
	T = newAppCnf()
}
