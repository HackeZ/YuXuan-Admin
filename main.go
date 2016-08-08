package main

import (
	_ "YuXuan-Admin/routers"
	"mime"
	"os"

	"YuXuan-Admin/models"
	"YuXuan-Admin/utils"

	"github.com/astaxie/beego"
)

// 初始化 DB
func initialize() {
	mime.AddExtensionType(".css", "text/css")

	// 判断初始化参数
	initArgs()

	// 连接数据库
	models.Connect()

	// 添加 StringsToJSON 方法
	beego.AddFuncMap("stringsToJson", utils.StringsToJSON)
}

func initArgs() {
	args := os.Args

	for _, v := range args {
		if v == "-syncdb" {
			models.Syncdb()
			os.Exit(0)
		}
	}
}

func main() {
	initialize()

	beego.Run()
}
