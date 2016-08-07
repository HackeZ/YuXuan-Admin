package main

import (
	_ "YuXuan-Admin/routers"
	"mime"

	"YuXuan-Admin/models"
	"YuXuan-Admin/utils"

	"github.com/astaxie/beego"
)

// 初始化 DB
func initialize() {
	mime.AddExtensionType(".css", "text/css")

	// 连接 MySQL 数据库
	models.Connect()

	// 添加 StringsToJSON 方法
	beego.AddFuncMap("stringsToJson", utils.StringsToJSON)
}

func main() {
	initialize()

	beego.Run()
}
