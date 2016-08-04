package controllers

import (
	"errors"
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
		// 默认认证类型 0 不认证 1 登录认证 2 实时认证
		userAuthType, _ := strconv.Atoi(beego.AppConfig.String("user_auth_type"))
		// 默认登录网关
		rbacAuthGateway := beego.AppConfig.String("rbac_auth_gateway")

		var accesslist map[string]bool
		if userAuthType > 0 {
			params := strings.Split(strings.ToLower(ctx.Request.RequestURI), "/")
		}
	}
}

// check user login
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
