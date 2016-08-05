package controllers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"

	m "YuXuan-Admin/models"
	"YuXuan-Admin/utils"
)

//AccessRegister check access and register user's nodes
func AccessRegister() {
	var Check = func(ctx *context.Context) {
		// 默认认证类型 0 不认证 1 登录认证 2 实时认证（不实现）
		userAuthType, _ := strconv.Atoi(beego.AppConfig.String("user_auth_type"))
		// 默认登录网关
		rbacAuthGateway := beego.AppConfig.String("rbac_auth_gateway")

		var accessList map[string]bool
		if userAuthType > 0 {
			// ctx.Requset.RequestURI --> http://localhost:8080/(public/login)/data/*
			// params --> ["public", "login"]
			params := strings.Split(strings.ToLower(ctx.Request.RequestURI), "/")

			if CheckAccess(params) {
				// verify success.
				uinfo := ctx.Input.Session("userinfo")
				if uinfo == nil {
					// have no userinfo, relogin.
					ctx.Redirect(302, rbacAuthGateway)
				}
				// Admin User No need to verify.
				adminuser := beego.AppConfig.String("rbac_admin_user")
				if adminuser == uinfo.(m.User).Username {
					return
				}

				// 登录验证
				if userAuthType == 1 {
					listBySession := ctx.Input.Session("accesslist")
					if listBySession == nil {
						accessList = listBySession.(map[string]bool)
					}
				}

				ret := AccessDecision(params, accessList)
				if !ret {
					ctx.Output.JSON(&map[string]interface{}{"status": false, "info": "权限不足"}, true, false)
				}
			}

		}
	}
	// 在访问 "/*" URI 时，寻找路由器之前(BeforeRouter)，先进行 Check 验证访问过滤器
	beego.InsertFilter("/*", beego.BeforeRouter, Check)
}

// CheckAccess Determine whether need to verify.
func CheckAccess(params []string) bool {
	if len(params) < 3 {
		return false
	}

	for _, nap := range strings.Split(beego.AppConfig.String("not_auth_package"), ",") {
		if params[1] == nap {
			// this URI need verify.
			return false
		}
	}
	return true
}

// AccessDecision To test wether permissions.
func AccessDecision(params []string, accesslist map[string]bool) bool {
	if CheckAccess(params) {
		s := fmt.Sprintf("%s/%s/%s", params[1], params[2], params[3])
		if len(accesslist) < 1 {
			return false
		}
		_, ok := accesslist[s]
		if ok != false {
			return true
		}
	} else {
		return true
	}
	return false
}

type AccessNode struct {
	Id        int64
	Name      string
	Childrens []*AccessNode
}

//Access permissions list
func GetAccessList(uid int64) (map[string]bool, error) {
	list, err := m.AccessList(uid)
	if err != nil {
		return nil, err
	}
	alist := make([]*AccessNode, 0)
	for _, l := range list {
		if l["Pid"].(int64) == 0 && l["Level"].(int64) == 1 {
			anode := new(AccessNode)
			anode.Id = l["Id"].(int64)
			anode.Name = l["Name"].(string)
			alist = append(alist, anode)
		}
	}
	for _, l := range list {
		if l["Level"].(int64) == 2 {
			for _, an := range alist {
				if an.Id == l["Pid"].(int64) {
					anode := new(AccessNode)
					anode.Id = l["Id"].(int64)
					anode.Name = l["Name"].(string)
					an.Childrens = append(an.Childrens, anode)
				}
			}
		}
	}
	for _, l := range list {
		if l["Level"].(int64) == 3 {
			for _, an := range alist {
				for _, an1 := range an.Childrens {
					if an1.Id == l["Pid"].(int64) {
						anode := new(AccessNode)
						anode.Id = l["Id"].(int64)
						anode.Name = l["Name"].(string)
						an1.Childrens = append(an1.Childrens, anode)
					}
				}

			}
		}
	}
	accesslist := make(map[string]bool)
	for _, v := range alist {
		for _, v1 := range v.Childrens {
			for _, v2 := range v1.Childrens {
				vname := strings.Split(v.Name, "/")
				v1name := strings.Split(v1.Name, "/")
				v2name := strings.Split(v2.Name, "/")
				str := fmt.Sprintf("%s/%s/%s", strings.ToLower(vname[0]), strings.ToLower(v1name[0]), strings.ToLower(v2name[0]))
				accesslist[str] = true
			}
		}
	}
	return accesslist, nil
}

// CheckLogin Check user login.
func CheckLogin(username, password string) (user m.User, err error) {
	user = m.GetUserByUsername(username)
	if user.Id == 0 {
		return user, errors.New("用户不存在")
	}
	if user.Password != utils.PassEncode(password, user.Salt) {
		return user, errors.New("密码错误")
	}
	return user, nil
}
