package Base

import (
	"fmt"
	//	"html/template"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/beego/i18n"

	"Hong/conf"
	"Hong/models/Farm"
	"Hong/models/model"
	"Hong/models/utils"
)

type NestPreparer interface {
	NestPrepare()
}

// baseRouter implemented global settings for all other routers.
type BaseRouter struct {
	beego.Controller
	i18n.Locale
	User    model.User
	IsLogin bool
}

// Prepare implemented Prepare method for baseRouter.
func (this *BaseRouter) Prepare() {
	//if conf.EnforceRedirect {
	//	// if the host not matching app settings then redirect to AppUrl
	//	if this.Ctx.Request.Host != conf.AppHost {
	//		this.Redirect(conf.AppUrl, 302)
	//		return
	//	}
	//}

	//// page start time
	this.Data["PageStartTime"] = time.Now()

	//// start session
	//this.StartSession()

	//// check flash redirect, if match url then end, else for redirect return
	//if match, redir := this.CheckFlashRedirect(this.Ctx.Request.RequestURI); redir {
	//	return
	//} else if match {
	//	this.EndFlashRedirect()
	//}

	//switch {
	//// save logined user if exist in session
	//case Farm.GetUserFromSession(&this.User, this.CruSession):
	//	this.IsLogin = true
	//// save logined user if exist in remember cookie
	//case Farm.LoginUserFromRememberCookie(&this.User, this.Ctx):
	//	this.IsLogin = true
	//}

	//if this.IsLogin {
	//	this.IsLogin = true
	//	this.Data["User"] = &this.User
	//	this.Data["IsLogin"] = this.IsLogin

	//	// if user forbided then do logout
	//	if this.User.IsForbid {
	//		Farm.LogoutUser(this.Ctx)
	//		this.FlashRedirect("/login", 302, "UserForbid")
	//		return
	//	}
	//}

	//// Setting properties.
	//this.Data["AppName"] = conf.AppName
	//this.Data["AppVer"] = conf.AppVersion
	//this.Data["AppUrl"] = conf.AppUrl
	//this.Data["AppLogo"] = conf.AppLogo
	//this.Data["AvatarURL"] = conf.AvatarURL
	//this.Data["IsProMode"] = conf.IsProMode
	//this.Data["SearchEnabled"] = conf.SearchEnabled
	//this.Data["NativeSearch"] = conf.NativeSearch
	//this.Data["SphinxEnabled"] = conf.SphinxEnabled

	//// Redirect to make URL clean.
	//if this.setLang() {
	//	i := strings.Index(this.Ctx.Request.RequestURI, "?")
	//	this.Redirect(this.Ctx.Request.RequestURI[:i], 302)
	//	return
	//}

	//// read flash message
	//beego.ReadFromRequest(&this.Controller)

	//// pass xsrf helper to template context
	//xsrfToken := this.XsrfToken()
	//this.Data["xsrf_token"] = xsrfToken
	//this.Data["xsrf_html"] = template.HTML(this.XsrfFormHtml())

	//// if method is GET then auto create a form once token
	//if this.Ctx.Request.Method == "GET" {
	//	//this.FormOnceCreate()
	//}

	//if app, ok := this.AppController.(NestPreparer); ok {
	//	app.NestPrepare()
	//}
	fmt.Println("Prepare() function")
}

// on router finished
func (this *BaseRouter) Finish() {
	fmt.Println("Finish() function")
}

func (this *BaseRouter) Get() {
	fmt.Println("Get() function")

}
func (this *BaseRouter) Post() {
	fmt.Println("Post() function")

}

func (this *BaseRouter) Delete() {
	fmt.Println("Delete() function")
}

func (this *BaseRouter) Put() {
	fmt.Println("Put() function")

}

func (this *BaseRouter) Head() {
	fmt.Println("Head() function")

}

func (this *BaseRouter) Patch() {
	fmt.Println("Patch() function")

}
func (this *BaseRouter) Options() {
	fmt.Println("Options() function")
}

func (this *BaseRouter) LoginUser(user *model.User, remember bool) string {
	fmt.Println("LoginUser() function")
	loginRedirect := strings.TrimSpace(this.Ctx.GetCookie("login_to"))
	if utils.IsMatchHost(loginRedirect) == false {
		loginRedirect = "/"
	} else {
		this.Ctx.SetCookie("login_to", "", -1, "/")
	}

	// login user
	Farm.LoginUser(user, this.Ctx, remember)

	this.setLangCookie(i18n.GetLangByIndex(user.Lang))

	return loginRedirect
}

// check if user not active then redirect
func (this *BaseRouter) CheckActiveRedirect(args ...interface{}) bool {
	fmt.Println("CheckActiveRedirect() function")
	var redirect_to string
	code := 302
	needActive := true
	for _, arg := range args {
		switch v := arg.(type) {
		case bool:
			needActive = v
		case string:
			// custom redirect url
			redirect_to = v
		case int:
			code = v
		}
	}
	if needActive {
		// check login
		if this.CheckLoginRedirect() {
			return true
		}

		// redirect to active page
		if !this.User.IsActive {
			this.FlashRedirect("/settings/profile", code, "NeedActive")
			return true
		}
	} else {
		// no need active
		if this.User.IsActive {
			if redirect_to == "" {
				redirect_to = "/"
			}
			this.Redirect(redirect_to, code)
			return true
		}
	}
	return false

}

// check if not login then redirect
func (this *BaseRouter) CheckLoginRedirect(args ...interface{}) bool {
	fmt.Println("CheckLoginRedirect() function")
	var redirect_to string
	code := 302
	needLogin := true
	for _, arg := range args {
		switch v := arg.(type) {
		case bool:
			needLogin = v
		case string:
			// custom redirect url
			redirect_to = v
		case int:
			// custom redirect url
			code = v
		}
	}

	// if need login then redirect
	if needLogin && !this.IsLogin {
		if len(redirect_to) == 0 {
			req := this.Ctx.Request
			scheme := "http"
			if req.TLS != nil {
				scheme += "s"
			}
			redirect_to = fmt.Sprintf("%s://%s%s", scheme, req.Host, req.RequestURI)
		}
		redirect_to = "/login?to=" + url.QueryEscape(redirect_to)
		this.Redirect(redirect_to, code)
		return true
	}

	// if not need login then redirect
	if !needLogin && this.IsLogin {
		if len(redirect_to) == 0 {
			redirect_to = "/"
		}
		this.Redirect(redirect_to, code)
		return true
	}
	return false
}

// read beego flash message
func (this *BaseRouter) FlashRead(key string) (string, bool) {
	fmt.Println("FlashRead() function")
	if data, ok := this.Data["flash"].(map[string]string); ok {
		value, ok := data[key]
		return value, ok
	}
	return "", false
}

// write beego flash message
func (this *BaseRouter) FlashWrite(key string, value string) {
	fmt.Println("FlashWrite() function")
	flash := beego.NewFlash()
	flash.Data[key] = value
	flash.Store(&this.Controller)
}

// check flash redirect, ensure browser redirect to uri and display flash message.
func (this *BaseRouter) CheckFlashRedirect(value string) (match bool, redirect bool) {
	fmt.Println("CheckFlashRedirect() function")
	v := this.GetSession("on_redirect")
	if params, ok := v.([]interface{}); ok {
		if len(params) != 5 {
			this.EndFlashRedirect()
			goto end
		}
		uri := utils.ToStr(params[0])
		code := 302
		if c, ok := params[1].(int); ok {
			if c/100 == 3 {
				code = c
			}
		}
		flag := utils.ToStr(params[2])
		flagVal := utils.ToStr(params[3])
		times := 0
		if v, ok := params[4].(int); ok {
			times = v
		}

		times += 1
		if times > 3 {
			// if max retry times reached then end
			this.EndFlashRedirect()
			goto end
		}

		// match uri or flash flag
		if uri == value || flag == value {
			match = true
		} else {
			// if no match then continue redirect
			this.FlashRedirect(uri, code, flag, flagVal, times)
			redirect = true
		}
	}
end:
	return match, redirect
}

// set flash redirect
func (this *BaseRouter) FlashRedirect(uri string, code int, flag string, args ...interface{}) {
	flagVal := "true"
	times := 0
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			flagVal = v
		case int:
			times = v
		}
	}

	if len(uri) == 0 || uri[0] != '/' {
		panic("flash reirect only support same host redirect")
	}

	params := []interface{}{uri, code, flag, flagVal, times}
	this.SetSession("on_redirect", params)

	this.FlashWrite(flag, flagVal)
	this.Redirect(uri, code)
}

// clear flash redirect
func (this *BaseRouter) EndFlashRedirect() {
	fmt.Println("EndFlashRedirect() function")
	this.DelSession("on_redirect")
}

func (this *BaseRouter) validForm(form interface{}, names ...string) (bool, map[string]*validation.ValidationError) {
	// parse request params to form ptr struct
	//utils.ParseForm(form, this.Input())
	fmt.Println("validForm() function")
	// Put data back in case users input invalid data for any section.
	name := reflect.ValueOf(form).Elem().Type().Name()
	if len(names) > 0 {
		name = names[0]
	}
	this.Data[name] = form

	errName := name + "Error"

	// check form once

	// Verify basic input.
	valid := validation.Validation{}
	if ok, _ := valid.Valid(form); !ok {
		errs := valid.ErrorMap()
		this.Data[errName] = &valid
		return false, errs
	}
	return true, nil
}

// valid form and put errors to tempalte context
func (this *BaseRouter) ValidForm(form interface{}, names ...string) bool {
	fmt.Println("ValidForm() function")
	valid, _ := this.validForm(form, names...)
	return valid
}

// add valid error to FormError
func (this *BaseRouter) SetFormError(form interface{}, fieldName, errMsg string, names ...string) {

	fmt.Println("SetFormError() function")
	name := reflect.ValueOf(form).Elem().Type().Name()
	if len(names) > 0 {
		name = names[0]
	}
	errName := name + "Error"
	//setsName := name + "Sets"

	if valid, ok := this.Data[errName].(*validation.Validation); ok {
		valid.SetError(fieldName, this.Tr(errMsg))
	}

	//if fSets, ok := this.Data[setsName].(*utils.FormSets); ok {
	//	fSets.SetError(fieldName, errMsg)
	//}
}

// check xsrf and show a friendly page
func (this *BaseRouter) CheckXsrfCookie() bool {
	fmt.Println("CheckXsrfCookie() function")
	return this.Controller.CheckXsrfCookie()
}

func (this *BaseRouter) IsAjax() bool {
	fmt.Println("IsAjax() function")
	return this.Ctx.Input.Header("X-Requested-With") == "XMLHttpRequest"
}

func (this *BaseRouter) setLangCookie(lang string) {
	fmt.Println("setLangCookie() function")
	this.Ctx.SetCookie("lang", lang, 60*60*24*365, "/", nil, nil, false)
}

// setLang sets site language version.
func (this *BaseRouter) setLang() bool {
	fmt.Println("setLang() function")
	isNeedRedir := false
	hasCookie := false

	// get all lang names from i18n
	langs := conf.Langs

	// 1. Check URL arguments.
	lang := this.GetString("lang")

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		lang = this.Ctx.GetCookie("lang")
		hasCookie = true
	} else {
		isNeedRedir = true
	}

	// Check again in case someone modify by purpose.
	if !i18n.IsExist(lang) {
		lang = ""
		isNeedRedir = false
		hasCookie = false
	}

	// 3. check if isLogin then use user setting
	if len(lang) == 0 && this.IsLogin {
		lang = i18n.GetLangByIndex(this.User.Lang)
	}

	// 4. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := this.Ctx.Input.Header("Accept-Language")
		if len(al) > 4 {
			al = al[:5] // Only compare first 5 letters.
			if i18n.IsExist(al) {
				lang = al
			}
		}
	}

	// 4. DefaucurLang language is English.
	if len(lang) == 0 {
		lang = "en-US"
		isNeedRedir = false
	}

	// Save language information in cookies.
	if !hasCookie {
		this.setLangCookie(lang)
	}

	// Set language properties.
	this.Data["Lang"] = lang
	this.Data["Langs"] = langs

	this.Lang = lang

	return isNeedRedir
}
