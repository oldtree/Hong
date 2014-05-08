// HongYu project HongYu.go
package main

import (
	"github.com/oldtree/Hong/conf"
	//	"Hong/models/Farm"
	//"Hong/models/FileManage"
	//"Hong/models/model"
	//	"Hong/models/utils"
	"github.com/oldtree/Hong/routers"
	"github.com/oldtree/Hong/routers/FarmPath"
	"fmt"
	"github.com/astaxie/beego"
)

func initall() {
	conf.LoadConfig()
	beego.ViewsPath = "views"
}

func main() {
	initall()
	beego.Info("AppPath:", beego.AppPath)

	beego.Info("AppPath:", beego.AppPath)

	if conf.IsProMode {
		beego.Info("Product mode enabled")
	} else {
		beego.Info("Develment mode enabled")
	}
	beego.Info(beego.AppName, conf.AppVersion, conf.AppUrl)

	fmt.Println("hello")
	world := new(routers.IntoWorld)
	beego.Router("/", world, "get:Get")
	login := new(FarmPath.LoginRouter)
	beego.Router("/login", login, "get:Get")

	beego.Run()

}
