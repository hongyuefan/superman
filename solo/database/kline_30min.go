package database

import (
	"errors"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type KLine_30Min struct {
	ID    int64  `orm:"column(id);auto"`
	Open  string `orm:"column(open);size(32)"`
	High  string `orm:"column(high);size(32)"`
	Low   string `orm:"column(low);size(32)"`
	Close string `orm:"column(close);size(32)"`
	Deal  string `orm:"column(deal);size(32)"`
	Time  string `orm:"column(time);size(32)"`
}

func (k *KLine_30Min) TableName() string {
	return "kline_30min"
}

func init() {

	orm.RegisterModel(new(KLine_30Min))

}

func SetKLine_30MinByTime(m *KLine_30Min) (num int64, err error) {

	o := orm.NewOrm()

	v := KLine_30Min{Time: m.Time}

	if err = o.Read(&v, "time"); err == nil {
		m.ID = v.ID
		return o.Update(m)
	}

	return AddKLine_30Min(m)
}

func AddKLine_30Min(m *KLine_30Min) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetKLine_30Min(id int64) (v *KLine_30Min, err error) {
	o := orm.NewOrm()
	v = &KLine_30Min{ID: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetKLine_30Mins(query map[string]string, fields []string, sortby []string, order []string, offset int64, limit int64) (ml []interface{}, err error) {

	o := orm.NewOrm()
	qs := o.QueryTable(new(KLine_30Min))
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

	var l []KLine_30Min
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
