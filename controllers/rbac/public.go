package rbac

import (
	"github.com/astaxie/beego"

	ctrl "YuXuan-Admin/controllers"
	m "YuXuanAPI/models"
)

// MainController ...
type MainController struct {
	CommonController
}

// Index 首页
func (c *MainController) Index() {
	userinfo := c.GetSession("userinfo")
	if userinfo == nil {
		c.Ctx.Redirect(302, beego.AppConfig.String("rbac_auth_gateway"))
	}
	// tree:=c.GetTree()
	if c.IsAjax() {
		c.Data["json"] = "Welcome to YuXuan-Admin Index"
		c.ServeJSON()
		return
	}
	groups := m.GroupList()
	c.Data["userinfo"] = userinfo
	c.Data["groups"] = groups
	// c.Data["tree"] = &tree
	if c.GetTemplatetype() != "easyui" {
		c.Layout = c.GetTemplatetype() + "/public/layout.tpl"
	}
	c.TplName = c.GetTemplatetype() + "/public/index.tpl"

}

// Login 登陆
func (c *MainController) Login() {
	isAjax := c.GetString("isajax")
	if isAjax == "1" {
		username := c.GetString("username")
		password := c.GetString("password")

		user, err := ctrl.CheckLogin(username, password)
		if err == nil {
			c.SetSession("userinfo", user)
			accessList, _ := ctrl.GetAccessList(user.Id)
			c.SetSession("accesslist", accessList)
			c.Resp(true, "登陆成功")
			return
		}
		c.Resp(false, err.Error())
		return

	}
	userinfo := c.GetSession("userinfo")
	if userinfo != nil {
		c.Ctx.Redirect(302, "/public/index")
	}
	c.TplName = c.GetTemplatetype() + "public/login.tpl"
}

// Logout 退出
func (c *MainController) Logout() {
	c.DelSession("userinfo")
	c.Ctx.Redirect(302, "/public/login")
}

// Changepwd 修改密码
func (c *MainController) Changepwd() {
	userinfo := c.GetSession("userinfo")
	if userinfo == nil {
		c.Ctx.Redirect(302, beego.AppConfig.String("rbac_auth_gateway"))
	}
	oldPassword := c.GetString("oldpassword")
	newPassword := c.GetString("newpassword")
	repeatPassword := c.GetString("repeatpassword")
	if newPassword != repeatPassword {
		c.Rsp(false, "两次输入密码不一致")
	}
	user, err := CheckLogin(userinfo.(m.User).Username, oldPassword)
	if err == nil {
		var u m.User
		u.Id = user.Id
		u.Password = newPassword
		id, err := m.UpdateUser(&u)
		if err == nil && id > 0 {
			c.Resp(true, "密码修改成功")
			return
		}
		c.Resp(false, err.Error())
		return
	}
	c.Resp(false, "密码有误")

}
