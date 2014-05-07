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

type IntoWorld struct {
	Base.BaseRouter
}

func (this *IntoWorld) Get() {
	fmt.Println("toTheDoor.Get")
	this.Data["IsLoginPage"] = false
	this.Data["title"] = "Welcome"
	this.TplNames = "world.html"
	this.Layout = "world.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["head"] = "labels/head.html"
	this.LayoutSections["bottom"] = "labels/bottom.html"
}

//func (this *LoginRouter) Put() {

//}
//func (this *LoginRouter) Put() {

//}
//func (this *LoginRouter) Put() {

//}
//func (this *LoginRouter) Put() {

//}
//func (this *LoginRouter) Put() {

//}
//func (this *LoginRouter) Put() {

//}
