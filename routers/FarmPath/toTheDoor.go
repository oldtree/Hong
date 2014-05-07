package FarmPath

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

type LoginRouter struct {
	Base.BaseRouter
}

func (this *LoginRouter) Get() {
	fmt.Println("toTheDoor.Get")
	this.Data["IsLoginPage"] = true
	this.Data["title"] = "hope"
	this.TplNames = "Myland/login.html"
	this.Layout = "Myland/login.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["head"] = "labels/head.html"
	this.LayoutSections["bottom"] = "labels/bottom.html"
	//this.LayoutSections["head"] = "labels/head.html"
	this.LayoutSections["body"] = "labels/body.html"
}

func (this *LoginRouter) Put() {

}
