package model

import (
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	//IsDel   soft_delete.DeletedAt `gorm:"softDelete:flag"` //软删除标识位 "gorm.io/plugin/soft_delete"
	UserID  uint   `gorm:"not null"`
	Name    string `gorm:"type:varchar(20) not null"`
	Phone   string `gorm:"type:varchar(11) not null"`
	Address string `gorm:"type:varchar(50) not null"`
}
