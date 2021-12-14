package doc

import (
	"gitlab.com/pakkaparn/dms-doc/user"
	"gorm.io/gorm"
)

type DocType string

const (
	File   DocType = "file"
	Folder         = "folder"
)

type Doc struct {
	gorm.Model
	Name   string  `json:"name" binding:"required"`
	Type   DocType `json:"type" binding:"required"`
	UserID uint    `json:"user_id" binding:"required"`
	User   user.User
	// ParentID *uint  `json:"child_id"`
	// Parent   *Doc   `gorm:"foreignKey:ParentID;references:ID"`
	// Children []*Doc `gorm:"foreignKey:ID;references:ParentID"`
}

func (Doc) TableName() string {
	return "docs"
}
