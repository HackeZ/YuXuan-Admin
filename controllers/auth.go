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

//check access and register user's nodes
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
