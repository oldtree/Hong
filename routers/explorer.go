package routers

import (
	//	"github.com/astaxie/beego"
	//	"strings"

	//	"Hong/conf"
	//	"Hong/models/Farm"
	//	"Hong/models/model"
	//	"Hong/models/utils"
	"Hong/routers/Base"
	"fmt"
)

type ExplorerWorld struct {
	Base.BaseRouter
}

func (this *ExplorerWorld) Get() {
	fmt.Println("ExplorerWorld.Get")
	this.Data["IsLoginPage"] = false
	this.Data["title"] = "Explorer this world"
	this.TplNames = "explorer.html"
	this.Layout = "explorer.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["head"] = "labels/head.html"
	this.LayoutSections["bottom"] = "labels/bottom.html"
}
