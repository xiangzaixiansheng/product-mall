package model

import (
	"product-mall/cache"
	"strconv"

	"gorm.io/gorm"
)

//商品信息
type Product struct {
	gorm.Model
	Name             string `gorm:"size:255;index"`
	CategoryID       uint   `gorm:"not null"`
	Title            string
	Info             string `gorm:"size:1000"`
	ImgPath          string
	Price            string
	DiscountPrice    string
	OnSale           bool `gorm:"default:false"`
	Num              int
	CreateUserID     int    //创建用户的Id
	CreateUserName   string // 创建用户的名字
	CreateUserAvatar string //创建用户的图像信息
}

// 获取点击数
func (product *Product) GetView() uint64 {
	countStr, _ := cache.GetInstance().Get(cache.GetViewKey(product.ID)).Result()
	count, _ := strconv.ParseUint(countStr, 10, 64)
	return count
}

//增加浏览量
func (product *Product) AddView() {
	// 增加视频点击数
	cache.GetInstance().Incr(cache.GetViewKey(product.ID))
	// 增加排行点击数
	cache.GetInstance().ZIncrBy(cache.RankKey, 1, strconv.Itoa(int(product.ID)))
}

// AddBookRank 图书排名
func (product *Product) AddBookRank() {
	//增加图书榜的点击量
	cache.GetInstance().ZIncrBy(cache.BookRank, 1, strconv.Itoa(int(product.ID)))
}

// AddCameraRank 相机排行榜
func (product *Product) AddCameraRank() {
	// 增加配件排行点击数
	cache.GetInstance().ZIncrBy(cache.CameraRank, 1, strconv.Itoa(int(product.ID)))
}
