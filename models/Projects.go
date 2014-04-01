package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	fmt.Println("init models file ")
	orm.RegisterModel(new(Project))
}

type Project struct {
	Id                         int64
	User                       *User     `orm:"rel(fk)"`
	Project_name               string    `orm:"size(256)"`
	Project_created            time.Time `orm:"auto_now_add;type(date)"`
	Project_update             time.Time `orm:"auto_now;type(date)"`
	Star                       int
	User_count                 int
	Project_history_Store_Path string `orm:"size(256)"`
}

func (u *Project) Insert() error {
	_, err := orm.NewOrm().Insert(u)
	return err
}

func (u *Project) Delete() error {
	_, err := orm.NewOrm().Delete(u)
	return err
}
func (u *Project) Read(fileds ...string) error {
	err := orm.NewOrm().Read(u, fileds...)
	return err
}

func (u *Project) Update(fileds ...string) error {
	_, err := orm.NewOrm().Update(u, fileds...)
	return err
}

func Projects() orm.QuerySeter {
	return orm.NewOrm().QueryTable("Project").OrderBy("-Id")
}
