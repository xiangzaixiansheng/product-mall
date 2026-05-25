package model

import "gorm.io/gorm"

//购物车
type Cart struct {
	gorm.Model
	UserID              uint //购物车创建的Id
	ProductID           uint `gorm:"not null"` //商品的ID
	ProductCreateUserID uint //商品的创建者的Id
	Num                 uint
	MaxNum              uint
	Check               bool
}
