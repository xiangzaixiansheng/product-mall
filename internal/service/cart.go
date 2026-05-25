package service

import (
	"product-mall/internal/dto"
	"product-mall/internal/model"
	"product-mall/pkg/db"
	"product-mall/pkg/e"
	"product-mall/pkg/pkg_logger"

	"strconv"
)

//购物车 解析form使用的
type CartService struct {
	CreateUserID uint `form:"create_user_id" json:"create_user_id"` //商品的创建人的id 后续应该升级为shopId
	Num          uint `from:"num" json:"num"`
}

//创建购物车
// 1、获取产品信息
// 2、根据产品信息 创建购物车信息 产品数量创建时时间等
// 3、已经在购物车的信息 添加的时候加一 不存在的时候创建
// id 商品id uid 用户id
func (service *CartService) Create(id string, uid uint) dto.Response {
	var product model.Product
	code := e.SUCCESS
	if err := db.GetDB().First(&product, id).Error; err != nil {
		//商品信息不存在也包含在里面
		pkg_logger.LogrusObj.Error("db error", "error", err)
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	cartId, _ := strconv.Atoi(id)
	var cart model.Cart
	db.GetDB().Where("user_id=? AND product_id=? AND product_create_user_id=?", uid, id, product.CreateUserID).First(&cart)
	if cart == (model.Cart{}) {
		//不存在购物车信息则创建
		cart = model.Cart{
			UserID:              uid,
			ProductID:           uint(cartId),
			ProductCreateUserID: uint(product.CreateUserID),
			Num:                 1,
			MaxNum:              10,
			Check:               false,
		}
		err := db.GetDB().Create(&cart).Error
		if err != nil {
			pkg_logger.LogrusObj.Error("db error", "error", err)
			code = e.ErrorDatabase
			return dto.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		return dto.Response{
			Status: code,
			Data:   dto.BuildCart(cart, product, service.CreateUserID),
			Msg:    e.GetMsg(code),
		}

	} else if cart.Num < cart.MaxNum {
		cart.Num++
		err := db.GetDB().Save(&cart).Error
		if err != nil {
			pkg_logger.LogrusObj.Error("db error", "error", err)
			return dto.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		return dto.Response{
			Status: 201,
			Msg:    "商品已经在购物车了，数量+1",
			Data:   dto.BuildCart(cart, product, service.CreateUserID),
		}

	} else {
		//超过最大的商品数量
		return dto.Response{
			Status: 202,
			Msg:    "超过购物车添加的最大上限",
		}
	}
}

//更新购物车信息
//id 购物车的id
func (service *CartService) Update(id string) dto.Response {
	var cart model.Cart
	code := e.SUCCESS
	err := db.GetDB().Where("id=?", id).Find(&cart).Error
	if err != nil {
		pkg_logger.LogrusObj.Error("db error", "error", err)
		code := e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	//更新数量
	cart.Num = service.Num
	err = db.GetDB().Save(&cart).Error
	if err != nil {
		pkg_logger.LogrusObj.Error("db error", "error", err)
		code := e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return dto.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}

}

//删除购物车信息
func (service *CartService) Delete(pid string, uid uint) dto.Response {
	var cart model.Cart
	code := e.SUCCESS
	err := db.GetDB().Where("user_id=? AND product_id=?", uid, pid).Error
	if err != nil {
		pkg_logger.LogrusObj.Error("db error", "error", err)
		code := e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	err = db.GetDB().Delete(&cart).Error
	if err != nil {
		pkg_logger.LogrusObj.Error("db error", "error", err)
		code := e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return dto.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}

}

//获取列表
//用户id
func (service *CartService) List(userId string) dto.Response {
	var carts []model.Cart
	code := e.SUCCESS
	err := db.GetDB().Where("user_id=?", userId).Find(&carts).Error

	if err != nil {
		pkg_logger.LogrusObj.Error("db error", "error", err)
		code := e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return dto.Response{
		Status: code,
		Data:   dto.BuildCarts(carts),
		Msg:    e.GetMsg(code),
	}

}
