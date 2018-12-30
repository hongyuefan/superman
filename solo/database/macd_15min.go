package database

import (
	"errors"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type MACD_15Min struct {
	ID    int64   `orm:"column(id);auto"`
	EMA12 float64 `orm:"column(ema12);"`
	EMA26 float64 `orm:"column(ema26);"`
	DIF   float64 `orm:"column(dif);"`
	DEA   float64 `orm:"column(dea);"`
	MACD  float64 `orm:"column(macd);"`
	Time  int64   `orm:"column(time);"`
}

func (m *MACD_15Min) TableName() string {
	return "macd_15min"
}

func init() {
	orm.RegisterModel(new(MACD_15Min))
}

func GetMACD_15Min(m *MACD_15Min, col ...string) (err error) {

	o := orm.NewOrm()

	return o.Read(m, col...)
}

func SetMACD_15Min(m *MACD_15Min) (num int64, err error) {

	o := orm.NewOrm()

	v := MACD_15Min{Time: m.Time}

	if err = o.Read(&v, "time"); err == nil {
		m.ID = v.ID
		return o.Update(m)
	}

	return AddMACD_15Min(m)
}

func AddMACD_15Min(m *MACD_15Min) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetMACD_15Min_Last(v []MACD_15Min, limit, offset int64) (int64, error) {
	return orm.NewOrm().QueryTable(new(MACD_15Min)).OrderBy("-id").Limit(limit, offset).All(&v)
}

func GetMACD_15Mins(query map[string]string, fields []string, sortby []string, order []string, offset int64, limit int64) (ml []interface{}, err error) {

	o := orm.NewOrm()
	qs := o.QueryTable(new(MACD_15Min))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []MACD_15Min
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}
