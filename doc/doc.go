package doc

import (
	"os/user"

	"gorm.io/gorm"
)

type DocType int64

const (
	File DocType = iota
	Folder
)

type Doc struct {
	gorm.Model
	Name string
	Type DocType
	Docs []Doc
	User user.User
}

func (Doc) TableName() string {
	return "docs"
}
