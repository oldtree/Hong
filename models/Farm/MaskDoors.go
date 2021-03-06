package Farm

import (
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego/context"
	"strings"
	// "time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/session"

	"Hong/conf"

	"Hong/models/model"
	"Hong/models/utils"
)

// CanRegistered checks if the username or e-mail is available.
func CanRegistered(userName string, email string) (bool, bool, error) {
	cond := orm.NewCondition()
	cond = cond.Or("UserName", userName).Or("Email", email)

	var maps []orm.Params
	o := orm.NewOrm()
	n, err := o.QueryTable("user").SetCond(cond).Values(&maps, "UserName", "Email")
	if err != nil {
		return false, false, err
	}

	e1 := true
	e2 := true

	if n > 0 {
		for _, m := range maps {
			if e1 && orm.ToStr(m["UserName"]) == userName {
				e1 = false
			}
			if e2 && orm.ToStr(m["Email"]) == email {
				e2 = false
			}
		}
	}

	return e1, e2, nil
}

// check if exist user by username or email
func HasUser(user *model.User, username string) bool {
	var err error
	qs := orm.NewOrm()
	if strings.IndexRune(username, '@') == -1 {
		user.UserName = username
		err = qs.Read(user, "UserName")
	} else {
		user.Email = username
		err = qs.Read(user, "Email")
	}
	if err == nil {
		return true
	}
	return false
}

// register create user
func RegisterUser(user *model.User, username, email, password string) error {
	// use random salt encode password
	salt := model.GetUserSalt()
	pwd := utils.EncodePassword(password, salt)

	user.UserName = strings.ToLower(username)
	user.Email = strings.ToLower(email)

	// save salt and encode password, use $ as split char
	user.Password = fmt.Sprintf("%s$%s", salt, pwd)

	// save md5 email value for gravatar
	user.GrEmail = utils.EncodeMd5(user.Email)

	// Use username as default nickname.
	user.NickName = user.UserName

	return user.Insert()
}

// set a new password to user
func SaveNewPassword(user *model.User, password string) error {
	salt := model.GetUserSalt()
	user.Password = fmt.Sprintf("%s$%s", salt, utils.EncodePassword(password, salt))
	return user.Update("Password", "Rands", "Updated")
}

// get login redirect url from cookie
func GetLoginRedirect(ctx *context.Context) string {
	loginRedirect := strings.TrimSpace(ctx.GetCookie("login_to"))
	if utils.IsMatchHost(loginRedirect) == false {
		loginRedirect = "/"
	} else {
		ctx.SetCookie("login_to", "", -1, "/")
	}
	return loginRedirect
}

// login user
func LoginUser(user *model.User, ctx *context.Context, remember bool) {
	// werid way of beego session regenerate id...
	ctx.Input.CruSession.SessionRelease(ctx.ResponseWriter)
	ctx.Input.CruSession = beego.GlobalSessions.SessionRegenerateId(ctx.ResponseWriter, ctx.Request)
	ctx.Input.CruSession.Set("auth_user_id", user.Id)

	if remember {
		WriteRememberCookie(user, ctx)
	}
}

func WriteRememberCookie(user *model.User, ctx *context.Context) {
	secret := utils.EncodeMd5(user.Rands + user.Password)
	days := 86400 * conf.LoginRememberDays
	ctx.SetCookie(conf.CookieUserName, user.UserName, days)
	ctx.SetSecureCookie(secret, conf.CookieRememberName, user.UserName, days)
}

func DeleteRememberCookie(ctx *context.Context) {
	ctx.SetCookie(conf.CookieUserName, "", -1)
	ctx.SetCookie(conf.CookieRememberName, "", -1)
}

func LoginUserFromRememberCookie(user *model.User, ctx *context.Context) (success bool) {
	userName := ctx.GetCookie(conf.CookieUserName)
	if len(userName) == 0 {
		return false
	}

	defer func() {
		if !success {
			DeleteRememberCookie(ctx)
		}
	}()

	user.UserName = userName
	if err := user.Read("UserName"); err != nil {
		return false
	}

	secret := utils.EncodeMd5(user.Rands + user.Password)
	value, _ := ctx.GetSecureCookie(secret, conf.CookieRememberName)
	if value != userName {
		return false
	}

	LoginUser(user, ctx, true)

	return true
}

// logout user
func LogoutUser(ctx *context.Context) {
	DeleteRememberCookie(ctx)
	ctx.Input.CruSession.Delete("auth_user_id")
	ctx.Input.CruSession.Flush()
	beego.GlobalSessions.SessionDestroy(ctx.ResponseWriter, ctx.Request)
}

func GetUserIdFromSession(sess session.SessionStore) int {
	if id, ok := sess.Get("auth_user_id").(int); ok && id > 0 {
		return id
	}
	return 0
}

// get user if key exist in session
func GetUserFromSession(user *model.User, sess session.SessionStore) bool {
	id := GetUserIdFromSession(sess)
	if id > 0 {
		u := model.User{Id: id}
		if u.Read() == nil {
			*user = u
			return true
		}
	}

	return false
}

// verify username/email and password
func VerifyUser(user *model.User, username, password string) (success bool) {
	// search user by username or email
	if HasUser(user, username) == false {
		return
	}

	if VerifyPassword(password, user.Password) {
		// success
		success = true

		// re-save discuz password
		if len(user.Password) == 39 {
			if err := SaveNewPassword(user, password); err != nil {
				beego.Error("SaveNewPassword err: ", err.Error())
			}
		}
	}
	return
}

// compare raw password and encoded password
func VerifyPassword(rawPwd, encodedPwd string) bool {

	// for discuz accounts
	if len(encodedPwd) == 39 {
		salt := encodedPwd[:6]
		encoded := encodedPwd[7:]
		return encoded == utils.EncodeMd5(utils.EncodeMd5(rawPwd)+salt)
	}

	// split
	var salt, encoded string
	if len(encodedPwd) > 11 {
		salt = encodedPwd[:10]
		encoded = encodedPwd[11:]
	}

	return utils.EncodePassword(rawPwd, salt) == encoded
}

// get user by erify code
func getVerifyUser(user *model.User, code string) bool {
	if len(code) <= utils.TimeLimitCodeLength {
		return false
	}

	// use tail hex username query user
	hexStr := code[utils.TimeLimitCodeLength:]
	if b, err := hex.DecodeString(hexStr); err == nil {
		user.UserName = string(b)
		if user.Read("UserName") == nil {
			return true
		}
	}

	return false
}

// verify active code when active account
func VerifyUserActiveCode(user *model.User, code string) bool {
	minutes := conf.ActiveCodeLives

	if getVerifyUser(user, code) {
		// time limit code
		prefix := code[:utils.TimeLimitCodeLength]
		data := utils.ToStr(user.Id) + user.Email + user.UserName + user.Password + user.Rands

		return utils.VerifyTimeLimitCode(data, minutes, prefix)
	}

	return false
}

// create a time limit code for user active
func CreateUserActiveCode(user *model.User, startInf interface{}) string {
	minutes := conf.ActiveCodeLives
	data := utils.ToStr(user.Id) + user.Email + user.UserName + user.Password + user.Rands
	code := utils.CreateTimeLimitCode(data, minutes, startInf)

	// add tail hex username
	code += hex.EncodeToString([]byte(user.UserName))
	return code
}

//// verify code when reset password
//func VerifyUserResetPwdCode(user *model.User, code string) bool {
//	minutes := conf.ResetPwdCodeLives

//	if getVerifyUser(user, code) {
//		// time limit code
//		prefix := code[:utils.TimeLimitCodeLength]
//		data := utils.ToStr(user.Id) + user.Email + user.UserName + user.Password + user.Rands + user.Updated.String()

//		return utils.VerifyTimeLimitCode(data, minutes, prefix)
//	}

//	return false
//}

//// create a time limit code for user reset password
//func CreateUserResetPwdCode(user *models.User, startInf interface{}) string {
//	minutes := conf.ResetPwdCodeLives
//	data := utils.ToStr(user.Id) + user.Email + user.UserName + user.Password + user.Rands + user.Updated.String()
//	code := utils.CreateTimeLimitCode(data, minutes, startInf)

//	// add tail hex username
//	code += hex.EncodeToString([]byte(user.UserName))
//	return code
//}
