package models

import (
	"errors"
	"log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

// 角色表
type Role struct {
	Id     int64
	Title  string `orm:"size(100)" form:"Title"  valid:"Required"`
	Name   string `orm:"size(100)" form:"Name"  valid:"Required"`
	Remark string `orm:"null;size(200)" form:"Remark" valid:"MaxSize(200)"`
	Status int    `orm:"default(2)" form:"Status" valid:"Range(1,2)"`
	//Node   []*Node `orm:"reverse(many)"`
	User []*User `orm:"reverse(many)"`
}

// TableName 定义角色表名称
func (r *Role) TableName() string {
	return beego.AppConfig.String("rbac_role_table")
}

// init 初始化Model
func init() {
	orm.RegisterModel(new(Role))
}

// checkRole 检查角色是否合法
func checkRole(g *Role) (err error) {
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

// GetRoleList get role list
func GetRoleList(page, pageSize int64, sort string) (roles []orm.Params, count int64) {
	o := orm.NewOrm()
	role := new(Role)
	qs := o.QueryTable(role)

	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * pageSize
	}

	qs.Limit(pageSize, offset).OrderBy(sort).Values(&roles)
	count, _ = qs.Count()
	return roles, count
}

// AddRole 添加角色
func AddRole(r *Role) (int64, error) {
	if err := checkRole(r); err != nil {
		return 0, err
	}

	o := orm.NewOrm()
	role := new(Role)
	role.Title = r.Title
	role.Name = r.Name
	role.Remark = r.Remark
	role.Status = r.Status

	id, err := o.Insert(role)
	return id, err
}

// UpdateRole 更新角色信息
func UpdateRole(r *Role) (int64, error) {
	if err := checkRole(r); err != nil {
		return 0, err
	}

	o := orm.NewOrm()
	role := make(orm.Params)
	if len(r.Title) > 0 {
		role["Title"] = r.Title
	}
	if len(r.Name) > 0 {
		role["Name"] = r.Name
	}
	if len(r.Remark) > 0 {
		role["Remark"] = r.Remark
	}
	if r.Status != 0 {
		role["Status"] = r.Status
	}
	if len(role) == 0 {
		return 0, errors.New("update field is empty")
	}

	var table Role
	num, err := o.QueryTable(table).Filter("Id", r.Id).Update(role)
	return num, err
}

// DelRoleById 删除角色
func DelRoleById(Id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Role{Id: Id})
	return status, err
}
