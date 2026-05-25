package dto

import (
	"product-mall/internal/model"
	"product-mall/pkg/db"
)

//购物车
type Cart struct {
	ID                  uint   `json:"id"`
	UserID              uint   `json:"user_id"`                //购物车创建的Id
	ProductID           uint   `json:"product_id"`             //商品的ID
	ProductCreateUserID uint   `json:"product_create_user_id"` //商品的创建者的Id
	Num                 uint   `json:"num"`                    //购物车中的数量
	MaxNum              uint   `json:"max_num"`                //最大可以加购的数量
	Check               bool   `json:"Check"`
	Price               string `json:"Price"`          //商品价格
	DiscountPrice       string `json:"discount_price"` //商品的促销价格
	Name                string `json:"name"`           //商品的名字
	ImgPath             string `json:"img_path"`       //商品的图片信息
	CreateAt            int64  `json:"create_at"`
}

// 拼接购物车信息
func BuildCart(item1 model.Cart, item2 model.Product, ProductCreateUserID uint) Cart {

	return Cart{
		ID:                  item1.ID,
		UserID:              item1.UserID,
		ProductID:           item1.ProductID,
		CreateAt:            item1.CreatedAt.Unix(),
		Num:                 item1.Num,
		MaxNum:              item1.MaxNum,
		Check:               false,
		Name:                item2.Name,
		ImgPath:             item2.ImgPath,
		Price:               item2.Price,
		DiscountPrice:       item2.DiscountPrice,
		ProductCreateUserID: ProductCreateUserID,
	}
}

func BuildCarts(items []model.Cart) (carts []Cart) {
	for _, cartInfo := range items {
		productInfo := model.Product{}
		productCreateUserID := cartInfo.ProductCreateUserID
		err := db.GetDB().First(&productInfo, cartInfo.ProductID, cartInfo.ProductCreateUserID).Error
		if err != nil {
			continue
		}
		cart := BuildCart(cartInfo, productInfo, productCreateUserID)
		carts = append(carts, cart)
	}
	//没有数据的情况下返回空数组
	if carts == nil {
		carts = make([]Cart, 0)
	}
	return carts
}
