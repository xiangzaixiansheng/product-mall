package model

import (
	"gorm.io/gorm"
)

//存放公告和邮件模板等
//后面通过id来找对应的模版信息

type Notice struct {
	gorm.Model
	Text string `gorm:"type:text"`
}
