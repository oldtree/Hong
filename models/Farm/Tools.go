package Farm

import (
	"github.com/astaxie/beego/orm"

	"fmt"

	"github.com/beego/i18n"

	"Hong/models/mail"

	"Hong/models/model"
	"Hong/models/utils"
)

func UserFollow(user *model.User, theUser *model.User) {
	if theUser.Read() == nil {
		var mutual bool
		tFollow := model.Follow{User: theUser, FollowUser: user}
		if err := tFollow.Read("User", "FollowUser"); err == nil {
			mutual = true
		}

		follow := model.Follow{User: user, FollowUser: theUser, Mutual: mutual}
		if err := follow.Insert(); err == nil && mutual {
			tFollow.Mutual = mutual
			tFollow.Update("Mutual")
		}

		if nums, err := user.FollowingUsers().Count(); err == nil {
			user.Following = int(nums)
			user.Update("Following")
		}

		if nums, err := theUser.FollowerUsers().Count(); err == nil {
			theUser.Followers = int(nums)
			theUser.Update("Followers")
		}
	}
}

func UserUnFollow(user *model.User, theUser *model.User) {
	num, _ := user.FollowingUsers().Filter("FollowUser", theUser.Id).Delete()
	if num > 0 {
		theUser.FollowingUsers().Filter("FollowUser", user.Id).Update(orm.Params{
			"Mutual": false,
		})

		if nums, err := user.FollowingUsers().Count(); err == nil {
			user.Following = int(nums)
			user.Update("Following")
		}

		if nums, err := theUser.FollowerUsers().Count(); err == nil {
			theUser.Followers = int(nums)
			theUser.Update("Followers")
		}
	}
}

// Send user register mail with active code
func SendRegisterMail(locale i18n.Locale, user *model.User) {
	code := CreateUserActiveCode(user, nil)

	subject := locale.Tr("mail.register_success_subject")

	data := mail.GetMailTmplData(locale.Lang, user)
	data["Code"] = code
	body := utils.RenderTemplate("mail/auth/register_success.html", data)

	msg := mail.NewMailMessage([]string{user.Email}, subject, body)
	msg.Info = fmt.Sprintf("UID: %d, send register mail", user.Id)

	// async send mail
	mail.SendAsync(msg)
}

// Send user reset password mail with verify code
func SendResetPwdMail(locale i18n.Locale, user *model.User) {
	//code := CreateUserResetPwdCode(user, nil)
	code := 0x00000
	subject := locale.Tr("mail.reset_password_subject")

	data := mail.GetMailTmplData(locale.Lang, user)
	data["Code"] = code
	body := utils.RenderTemplate("mail/auth/reset_password.html", data)

	msg := mail.NewMailMessage([]string{user.Email}, subject, body)
	msg.Info = fmt.Sprintf("UID: %d, send reset password mail", user.Id)

	// async send mail
	mail.SendAsync(msg)
}

// Send email verify active email.
func SendActiveMail(locale i18n.Locale, user *model.User) {
	code := CreateUserActiveCode(user, nil)

	subject := locale.Tr("mail.verify_your_email_subject")

	data := mail.GetMailTmplData(locale.Lang, user)
	data["Code"] = code
	body := utils.RenderTemplate("mail/auth/active_email.html", data)

	msg := mail.NewMailMessage([]string{user.Email}, subject, body)
	msg.Info = fmt.Sprintf("UID: %d, send email verify mail", user.Id)

	// async send mail
	mail.SendAsync(msg)
}
