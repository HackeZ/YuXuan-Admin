package rbac

import (
	ctrl "YuXuan-Admin/controllers"

	"github.com/astaxie/beego"
)

// CommonController rbac common controllers
type CommonController struct {
	beego.Controller
	Templatetype string // ui template type
}

// Resp Response request Status and Info.
func (c *CommonController) Resp(status bool, str string) {
	c.Data["json"] = &map[string]interface{}{"status": status, "info": str}
	c.ServeJSON()
}

// GetTemplatetype set template theme.
func (c *CommonController) GetTemplatetype() string {
	templatetype := beego.AppConfig.String("template_type")
	if templatetype == "" {
		templatetype = "easyui"
	}
	return templatetype
}

func init() {
	// 验证权限
	ctrl.AccessRegister()
}
