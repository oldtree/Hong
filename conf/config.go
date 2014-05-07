// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package conf

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "mysql"
	"os"
	"path"
)

/*



   AppName

   应用名称，默认是 beego。通过bee new创建的是创建的项目名。

   AppPath

   当前应用的路径，默认会通过设置os.Args[0]获得执行的命令的第一个参数，所以你在使用 supervisor 管理进程的时候记得采用全路径启动。

   AppConfigPath

   配置文件所在的路径，默认是应用程序对应的目录下的 conf/app.conf，用户可以修改该值配置自己的配置文件。

   HttpAddr

   应用监听地址，默认为空，监听所有的网卡 IP。

   HttpPort

   应用监听端口，默认为 8080。

   HttpTLS

   是否启用 HTTPS，默认是关闭。

   HttpCertFile

   开启 HTTPS 之后，certfile 的路径。

   HttpKeyFile

   开启 HTTPS 之后，keyfile 的路径。

   HttpServerTimeOut

   设置 HTTP 的超时时间，默认是 0，不超时。

   RunMode

   应用的模式，默认是 dev，为开发模式，在开发模式下出错会提示友好的出错页面，如前面错误描述中所述。

   AutoRender

   是否模板自动渲染，默认值为 true，对于 API 类型的应用，应用需要把该选项设置为 false，不需要渲染模板。

   RecoverPanic

   是否异常恢复，默认值为 true，即当应用出现异常的情况，通过 recover 恢复回来，而不会导致应用异常退出。

   ViewsPath

   模板路径，默认值是 views。

   SessionOn

   session 是否开启，默认是 false。

   SessionProvider

   session 的引擎，默认是 memory。

   SessionName

   存在客户端的 cookie 名称，默认值是 beegosessionID。

   SessionGCMaxLifetime

   session 过期时间，默认值是 3600 秒。

   SessionSavePath

   session 保存路径，默认是空。

   SessionHashFunc

   sessionID 生成函数，默认是 sha1。

   SessionHashKey

   session hash 的 key。

   SessionCookieLifeTime

   session 默认存在客户端的 cookie 的时间，默认值是 3600 秒。

   UseFcgi

   是否启用 fastcgi，默认是 false。

   MaxMemory

   文件上传默认内存缓存大小，默认值是 1 << 26(64M)。

   EnableGzip

   是否开启 gzip 支持，默认为 false 不支持 gzip，一旦开启了 gzip，那么在模板输出的内容会进行 gzip 或者 zlib 压缩，根据用户的 Accept-Encoding 来判断。

   DirectoryIndex

   是否开启静态目录的列表显示，默认不显示目录，返回 403 错误。

   BeegoServerName

   beego 服务器默认在请求的时候输出 server 为 beego。

   EnableAdmin

   是否开启进程内监控模块，默认关闭。

   AdminHttpAddr

   监控程序监听的地址，默认值是 localhost。

   AdminHttpPort

   监控程序监听的端口，默认值是 8088。

   TemplateLeft

   模板左标签，默认值是{{。

   TemplateRight

   模板右标签，默认值是}}。

   ErrorsShow

   是否显示错误，默认显示错误信息。

   XSRFKEY

   XSRF 的 key 信息，默认值是 beegoxsrf。

   XSRFExpire

   XSRF 过期时间，默认值是 0。



*/
var (
	AppName         string
	AppPath         string
	AppVersion      string
	AppHost         string
	AppUrl          string
	AppLogo         string
	EnforceRedirect bool
	AvatarURL       string
	SecretKey       string
	IsProMode       bool
	//MailUser            string
	//MailFrom            string
	ActiveCodeLives   int
	ResetPwdCodeLives int

	DateFormat          string
	DateTimeFormat      string
	DateTimeShortFormat string
	TimeZone            string
	RealtimeRenderMD    bool
	ImageSizeSmall      int
	ImageSizeMiddle     int
	ImageLinkAlphabets  []byte
	ImageXSend          bool
	ImageXSendHeader    string
	Langs               []string

	LoginRememberDays int
	LoginMaxRetries   int
	LoginFailedBlocks int

	CookieRememberName string
	CookieUserName     string

	// search
	SearchEnabled bool
	NativeSearch  bool

	// sphinx search setting
	SphinxEnabled bool
	SphinxHost    string
	SphinxIndex   string
	SphinxMaxConn int

	// mail setting
	MailUser     string
	MailFrom     string
	MailHost     string
	MailAuthUser string
	MailAuthPass string
)

var (
	// OAuth
	GithubClientId     string
	GithubClientSecret string
	GoogleClientId     string
	GoogleClientSecret string
	WeiboClientId      string
	WeiboClientSecret  string
	QQClientId         string
	QQClientSecret     string
)

const (
	AppConfigPath = "conf/app.ini"
)
const (
	LangEnUS = iota
	LangZhCN
)

var (
	Cfg *goconfig.ConfigFile
)

func loadConfigVar() {
	AppName = beego.AppName

	AppHost = Cfg.MustValue("app", "app_host")
	AppUrl = Cfg.MustValue("app", "app_url")
	AppLogo = Cfg.MustValue("app", "app_logo")
	AppVersion = Cfg.MustValue("app", "app_ver")

	ActiveCodeLives = Cfg.MustInt("app", "acitve_code_live_hours")

	if ActiveCodeLives <= 0 {
		ActiveCodeLives = 12
	}
	ResetPwdCodeLives = Cfg.MustInt("app", "resetpwd_code_live_hours")
	if ResetPwdCodeLives <= 0 {
		ResetPwdCodeLives = 12
	}

	LoginRememberDays = Cfg.MustInt("app", "login_remember_days")

	LoginMaxRetries = Cfg.MustInt("app", "login_max_retries")
	if LoginMaxRetries <= 0 {
		LoginMaxRetries = 3
	}

	LoginFailedBlocks = Cfg.MustInt("app", "login_failed_blocks")
	if LoginFailedBlocks <= 0 {
		LoginFailedBlocks = 5
	}

}

func LoadConfig() {

	if fh, _ := os.OpenFile(AppConfigPath, os.O_RDONLY|os.O_CREATE, 0777); fh != nil {
		fh.Close()
	}

	var err error

	Cfg, err = goconfig.LoadConfigFile(AppConfigPath)

	Cfg.BlockMode = false

	if err != nil {
		panic("Fail to load configuration file: " + err.Error())
	}
	beego.AppName = Cfg.MustValue("app", "app_name")
	beego.HttpPort = Cfg.MustInt("app", "http_port")
	beego.RunMode = Cfg.MustValue("app", "run_mode")

	IsProMode = (beego.RunMode == "pro")
	if IsProMode {
		beego.SetLevel(beego.LevelInfo)
	}

	//session
	beego.SessionOn = true
	beego.SessionProvider = Cfg.MustValue("session", "session_provider")
	beego.SessionSavePath = Cfg.MustValue("session", "session_path")
	beego.SessionName = Cfg.MustValue("session", "session_name")

	//database
	driverName := Cfg.MustValue("orm", "driver_name")
	dataSource := Cfg.MustValue("orm", "data_source")
	maxIdle := Cfg.MustInt("orm", "max_idle_conn")
	maxOpen := Cfg.MustInt("orm", "max_open_conn")

	//set default database
	if _, err := os.Open(dataSource); err != nil && os.IsNotExist(err) {
		os.MkdirAll(path.Dir(dataSource), os.ModePerm)
		os.Create(dataSource)
	}
	orm.RegisterDataBase("default", driverName, dataSource, maxIdle, maxOpen)
	orm.RunCommand()
	orm.Debug = true
	err = orm.RunSyncdb("default", false, true)

	orm.Debug = Cfg.MustBool("orm", "debug_log")

	loadConfigVar()

}
func ShowInit() {
	fmt.Println("Done")
}
