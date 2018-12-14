package config

// app config
type AppCnf struct {
	LogPath string
	CnfPath string
	SqlCon  string
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

	return err
}

func newAppCnf() *AppCnf {
	return &AppCnf{}
}

func init() {
	T = newAppCnf()
}
