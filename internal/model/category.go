package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name     string `gorm:"size:100;not null;index"`
	ParentID uint   `gorm:"default:0;index"`
	Level    uint   `gorm:"default:1"`
	Sort     int    `gorm:"default:0"`
}
