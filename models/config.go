package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type Config struct {
	Id        int64  `orm:"column(id);auto"`
	Kafka     string `orm:"column(kafka);size(256);"`
	InfluxDB  string `orm:"column(influxDB);size(256);"`
	ExChanges string `orm:"column(exchanges);size(256);"` //okex,huobi
	Symbols   string `orm:"column(symbols);size(256)"`
	KlineTime string `orm:"column(kiline_time);size(256)"`
	Depth     string `orm:"column(depth);size(256)"`
	Path      string `orm:"column(store_path);size(256)"`
}

func (t *Config) TableName() string {
	return "explat_config"
}

func init() {
	orm.RegisterModel(new(Config))
}

func AddConfig(m *Config) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetConfig() (v *Config, err error) {
	o := orm.NewOrm()
	v = &Config{Id: 1}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetConfigById(id int64) (v *Config, err error) {
	o := orm.NewOrm()
	v = &Config{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func UpdateConfigById(m *Config, cols ...string) (err error) {
	o := orm.NewOrm()
	v := Config{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m, cols...); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func DeleteConfig(id int64) (err error) {
	o := orm.NewOrm()
	v := Config{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Config{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
