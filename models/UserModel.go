package models

import (
	"errors"
	"log"
	"time"

	"YuXuan-Admin/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

// User 用户表
type User struct {
	Id            int64
	Username      string    `orm:"unique;size(32)" form:"Username"  valid:"Required;MaxSize(20);MinSize(6)"`
	Password      string    `orm:"size(32)" form:"Password" valid:"Required;MaxSize(20);MinSize(6)"`
	Repassword    string    `orm:"-" form:"Repassword" valid:"Required"`
	Salt          string    `orm:"size(32)" form:"Salt" valid:"Required;MaxSize(20);MinSize(6)"`
	Nickname      string    `orm:"unique;size(32)" form:"Nickname" valid:"Required;MaxSize(20);MinSize(2)"`
	Email         string    `orm:"size(32)" form:"Email" valid:"Email"`
	Remark        string    `orm:"null;size(200)" form:"Remark" valid:"MaxSize(200)"`
	Status        int       `orm:"default(2)" form:"Status" valid:"Range(1,2)"`
	Lastlogintime time.Time `orm:"null;type(datetime)" form:"-"`
	Createtime    time.Time `orm:"type(datetime);auto_now_add" `
	Role          []*Role   `orm:"rel(m2m)"`
}

// TableName User Table Name.
func (u *User) TableName() string {
	return beego.AppConfig.String("rbac_user_table")
}

func init() {
	orm.RegisterModel(new(User))
}

// Valid 密码验证
func (u *User) Valid(v *validation.Validation) {
	if u.Password != u.Repassword {
		v.SetError("Repassword", "两次输入的密码不一样")
	}
}

//验证用户信息
func checkUser(u *User) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(&u)
	if !b {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			return errors.New(err.Message)
		}
	}
	return nil
}

/************************************************************/

// Getuserlist get user list
func Getuserlist(page int64, pageSize int64, sort string) (users []orm.Params, count int64) {
	o := orm.NewOrm()
	user := new(User)
	qs := o.QueryTable(user)
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * pageSize
	}
	qs.Limit(pageSize, offset).OrderBy(sort).Values(&users)
	count, _ = qs.Count()
	return users, count
}

// AddUser 添加用户
func AddUser(u *User) (int64, error) {
	if err := checkUser(u); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	user := new(User)
	user.Username = u.Username

	// Get New Salt and Handle Password.
	var saltLength int64
	if beego.AppConfig.String("salt_length") == "" {
		saltLength = 8
	} else {
		saltLength, _ = utils.Atoi64(beego.AppConfig.String("salt_length"))
	}
	userSalt := utils.GetNewSalt(saltLength)
	user.Salt = userSalt
	user.Password = utils.PassEncode(u.Password, user.Salt)

	user.Nickname = u.Nickname
	user.Email = u.Email
	user.Remark = u.Remark
	user.Status = u.Status

	id, err := o.Insert(user)
	return id, err
}

// UpdateUser 更新用户
func UpdateUser(u *User) (int64, error) {
	if err := checkUser(u); err != nil {
		return 0, err
	}

	o := orm.NewOrm()
	user := make(orm.Params)
	if len(u.Username) > 0 {
		user["Username"] = u.Username
	}
	if len(u.Nickname) > 0 {
		user["Nickname"] = u.Nickname
	}
	if len(u.Email) > 0 {
		user["Email"] = u.Email
	}
	if len(u.Remark) > 0 {
		user["Remark"] = u.Remark
	}
	if len(u.Password) > 0 {
		var saltLength int64
		if beego.AppConfig.String("salt_length") == "" {
			saltLength = 8
		} else {
			saltLength, _ = utils.Atoi64(beego.AppConfig.String("salt_length"))
		}
		userSalt := utils.GetNewSalt(saltLength)
		user["Salt"] = userSalt
		user["Password"] = utils.PassEncode(u.Password, userSalt)
	}
	if u.Status != 0 {
		user["Status"] = u.Status
	}
	if len(user) == 0 {
		return 0, errors.New("update field is empty")
	}

	var table User
	num, err := o.QueryTable(table).Filter("Id", u.Id).Update(user)
	return num, err
}

// DelUserById 删除用户
func DelUserById(Id int64) (int64, error) {
	o := orm.NewOrm()

	status, err := o.Delete(&User{Id: Id})
	return status, err
}

// GetUserByUsername 获取用户信息
func GetUserByUsername(username string) (user User) {
	user = User{Username: username}
	o := orm.NewOrm()
	log.Println("user -->", user)

	o.Read(&user, "Username")
	return user
}

// GetSaltById 获取用户盐
func GetSaltById(Id int64) (string, error) {
	user := User{Id: Id}
	o := orm.NewOrm()

	err := o.Read(&user)
	return user.Salt, err
}
