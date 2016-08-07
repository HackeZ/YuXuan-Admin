package models

//分组表
type Group struct {
	Id     int64
	Name   string  `orm:"size(100)" form:"Name"  valid:"Required"`
	Title  string  `orm:"size(100)" form:"Title"  valid:"Required"`
	Status int     `orm:"default(2)" form:"Status" valid:"Range(1,2)"`
	Sort   int     `orm:"default(1)" form:"Sort" valid:"Numeric"`
	Nodes  []*Node `orm:"reverse(many)"`
}
