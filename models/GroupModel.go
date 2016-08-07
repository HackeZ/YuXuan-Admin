package models

import (
	"errors"
	"log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

// Group 分组表
type Group struct {
	Id     int64
	Name   string  `orm:"size(100)" form:"Name"  valid:"Required"`
	Title  string  `orm:"size(100)" form:"Title"  valid:"Required"`
	Status int     `orm:"default(2)" form:"Status" valid:"Range(1,2)"`
	Sort   int     `orm:"default(1)" form:"Sort" valid:"Numeric"`
	Nodes  []*Node `orm:"reverse(many)"`
}

// TableName Group Table Name
func (g *Group) TableName() string {
	return beego.AppConfig.String("rbac_group_table")
}

func init() {
	orm.RegisterModel(new(Group))
}

// checkGroup  检查Group合法性
func checkGroup(g *Group) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(&g)

	if !b {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			return errors.New(err.Message)
		}
	}
	return nil
}

// GetGrouplist get group list
func GetGrouplist(page int64, pageSize int64, sort string) (groups []orm.Params, count int64) {
	o := orm.NewOrm()
	group := new(Group)
	qs := o.QueryTable(group)
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * pageSize
	}
	qs.Limit(pageSize, offset).OrderBy(sort).Values(&groups)
	count, _ = qs.Count()
	return groups, count
}

// AddGroup 添加分组
func AddGroup(g *Group) (int64, error) {
	if err := checkGroup(g); err != nil {
		return 0, err
	}

	o := orm.NewOrm()
	group := new(Group)
	group.Name = g.Name
	group.Title = g.Title
	group.Sort = g.Sort
	group.Status = g.Status

	id, err := o.Insert(group)
	return id, err
}

// UpdateGroup 更新分组信息
func UpdateGroup(g *Group) (int64, error) {
	if err := checkGroup(g); err != nil {
		return 0, err
	}

	o := orm.NewOrm()
	group := make(orm.Params)
	if len(g.Name) > 0 {
		group["Name"] = g.Name
	}
	if len(g.Title) > 0 {
		group["Title"] = g.Title
	}
	if g.Status != 0 {
		group["Status"] = g.Status
	}
	if g.Sort != 0 {
		group["Sort"] = g.Sort
	}
	if len(group) == 0 {
		return 0, errors.New("update field is empty")
	}

	var table Group
	num, err := o.QueryTable(table).Filter("Id", g.Id).Update(group)
	return num, err
}

// DelGroupById 删除分组
func DelGroupById(Id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Group{Id: Id})
	return status, err
}

// GroupList 分组列表
func GroupList() (groups []orm.Params) {
	o := orm.NewOrm()
	group := new(Group)
	qs := o.QueryTable(group)

	qs.Values(&groups, "id", "title")
	return groups
}
