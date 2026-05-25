package model

import "gorm.io/gorm"

const (
	OrderStatusPending   = 0 // 待支付
	OrderStatusPaid      = 1 // 已支付
	OrderStatusShipped   = 2 // 已发货
	OrderStatusCompleted = 3 // 已完成
	OrderStatusCancelled = 4 // 已取消
)

type Order struct {
	gorm.Model
	OrderNo    string  `gorm:"uniqueIndex;size:64;not null"`
	UserID     uint    `gorm:"index;not null"`
	AddressID  uint    `gorm:"not null"`
	TotalPrice float64 `gorm:"not null"`
	Status     int     `gorm:"default:0;index"`
	Remark     string  `gorm:"size:500"`
}

type OrderItem struct {
	gorm.Model
	OrderID     uint    `gorm:"index;not null"`
	ProductID   uint    `gorm:"not null"`
	ProductName string  `gorm:"size:255"`
	Num         uint    `gorm:"not null"`
	Price       float64 `gorm:"not null"`
}
