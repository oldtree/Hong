package utils

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"

	"Hong/conf"
)

// get HTML i18n string
func i18nHTML(lang, format string, args ...interface{}) template.HTML {
	return template.HTML(i18n.Tr(lang, format, args...))
}

func boolicon(b bool) (s template.HTML) {
	if b {
		s = `<i style="color:green;" class="icon-check""></i>`
	} else {
		s = `<i class="icon-check-empty""></i>`
	}
	return
}

func date(t time.Time) string {
	return beego.Date(t, conf.DateFormat)
}

func datetime(t time.Time) string {
	return beego.Date(t, conf.DateTimeFormat)
}

func datetimes(t time.Time) string {
	return beego.Date(t, conf.DateTimeShortFormat)
}

func loadtimes(t time.Time) int {
	return int(time.Since(t).Nanoseconds() / 1e6)
}

func sum(base interface{}, value interface{}, params ...interface{}) (s string) {
	switch v := base.(type) {
	case string:
		s = v + ToStr(value)
		for _, p := range params {
			s += ToStr(p)
		}
	}
	return s
}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func timesince(lang string, t time.Time) string {
	seconds := int(time.Since(t).Seconds())
	switch {
	case seconds < 60:
		return i18n.Tr(lang, "seconds_ago", seconds)
	case seconds < 60*60:
		return i18n.Tr(lang, "minutes_ago", seconds/60)
	case seconds < 60*60*24:
		return i18n.Tr(lang, "hours_ago", seconds/(60*60))
	case seconds < 60*60*24*100:
		return i18n.Tr(lang, "days_ago", seconds/(60*60*24))
	default:
		return beego.Date(t, conf.DateFormat)
	}
}

// create an login url with specify redirect to param
func loginto(uris ...string) template.HTMLAttr {
	var uri string
	if len(uris) > 0 {
		uri = uris[0]
	}
	to := fmt.Sprintf("%slogin", conf.AppUrl)
	if len(uri) > 0 {
		to += "?to=" + url.QueryEscape(uri)
	}
	return template.HTMLAttr(to)
}

func RenderTemplate(TplNames string, Data map[interface{}]interface{}) string {
	if beego.RunMode == "dev" {
		beego.BuildTemplate(beego.ViewsPath)
	}

	ibytes := bytes.NewBufferString("")
	if _, ok := beego.BeeTemplates[TplNames]; !ok {
		panic("can't find templatefile in the path:" + TplNames)
	}
	err := beego.BeeTemplates[TplNames].ExecuteTemplate(ibytes, TplNames, Data)
	if err != nil {
		beego.Trace("template Execute err:", err)
	}
	icontent, _ := ioutil.ReadAll(ibytes)
	return string(icontent)
}

func init() {
	// Register template functions.
	fmt.Println("template.init function")
	beego.AddFuncMap("i18n", i18nHTML)
	beego.AddFuncMap("boolicon", boolicon)
	beego.AddFuncMap("date", date)
	beego.AddFuncMap("datetime", datetime)
	beego.AddFuncMap("datetimes", datetimes)
	beego.AddFuncMap("dict", dict)
	beego.AddFuncMap("timesince", timesince)
	beego.AddFuncMap("loadtimes", loadtimes)
	beego.AddFuncMap("sum", sum)
	beego.AddFuncMap("loginto", loginto)
}
