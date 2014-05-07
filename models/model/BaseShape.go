package model

import (
	"fmt"
	"time"
)

type Zero interface {
	GetAuthorList(projectname string)
	GetProjectName(id int)
}

type Things struct {
	Id               int
	User             *User     `orm:"rel(fk)"`
	Uri              string    `orm:"size(60);unqiue"`
	Title            string    `orm:"size(60)"`
	Content          string    `orm:"type(text)"`
	ContentCache     string    `orm:"type(text)"`
	TitleZhCn        string    `orm:"size(60)"`
	ContentZhCn      string    `orm:"type(text)"`
	ContentCacheZhCn string    `orm:"type(text)"`
	LastAuthor       *User     `orm:"rel(fk);null"`
	IsPublish        bool      `orm:"index"`
	Created          time.Time `orm:"auto_now_add"`
	Updated          time.Time `orm:"auto_now"`
}

func init() {
	fmt.Println("init base shape ")
}

func ShowBae() {
	fmt.Println("两天之内完成 ")
}
