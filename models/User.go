package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	fmt.Println("init models file ")
	orm.RegisterModel(new(User))
}

type User struct {
	Id            int64
	User_name     string    `orm:"size(128)"`
	User_email    string    `orm:"size(128);unique"`
	User_address  string    `orm:"size(128)"`
	User_password string    `orm:"size(128)"`
	User_created  time.Time `orm:"auto_now_add;type(date)"`
	User_update   time.Time `orm:"auto_now;type(date)"`
	User_belong   string    `orm:"size(128)"`
	//Project       *Project  `orm:"null;rel(one);on_delete(set_null)"`
	Project []*Project `orm:"reverse(many)"`
}

func (u *User) Insert() error {
	_, err := orm.NewOrm().Insert(u)
	return err
}

func (u *User) Delete() error {
	_, err := orm.NewOrm().Delete(u)
	return err
}
func (u *User) Read(fileds ...string) error {
	err := orm.NewOrm().Read(u, fileds...)
	return err
}

func (u *User) Update(fileds ...string) error {
	_, err := orm.NewOrm().Update(u, fileds...)
	return err
}

func Users() orm.QuerySeter {
	return orm.NewOrm().QueryTable("User").OrderBy("-Id")
}
