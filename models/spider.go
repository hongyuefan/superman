package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type Spider struct {
	Id        int64  `orm:"column(id);auto"`
	ExName    string `orm:"column(exchange);size(32)"`
	WsUrl     string `orm:"column(ws_url);size(256)"`
	RpcUrl    string `orm:"column(rpc_url);size(256)"`
	Symbols   string `orm:"column(symbols);size(256)"`
	KlineTime string `orm:"column(kiline_time);size(256)"`
	Depth     string `orm:"column(depth);size(256)"`
	HeartBeat int64  `orm:"column(heartbeat);"`
}

func (t *Spider) TableName() string {
	return "explat_spider"
}

func init() {
	orm.RegisterModel(new(Spider))
}

func AddSpider(m *Spider) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetSpiderByName(exName string) (v *Spider, err error) {
	o := orm.NewOrm()
	v = &Spider{ExName: exName}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func UpdateSpiderByName(m *Spider, cols ...string) (err error) {
	o := orm.NewOrm()
	v := Spider{ExName: m.ExName}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m, cols...); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func DeleteSpider(exName string) (err error) {
	o := orm.NewOrm()
	v := Spider{ExName: exName}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Spider{ExName: exName}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
