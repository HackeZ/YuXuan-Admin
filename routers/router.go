package routers

import (
	"YuXuan-Admin/controllers"

	"github.com/astaxie/beego"
)

func init() {
	AdminRouter()
	beego.Router("/", &controllers.MainController{})
}
