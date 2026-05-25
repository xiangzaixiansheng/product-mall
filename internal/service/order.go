package service

import (
	"fmt"
	"product-mall/cache"
	"product-mall/internal/dto"
	"product-mall/internal/model"
	"product-mall/pkg/db"
	"product-mall/pkg/e"
	"product-mall/pkg/pkg_logger"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateOrderService struct {
	AddressID uint   `form:"address_id" json:"address_id" binding:"required"`
	Remark    string `form:"remark" json:"remark"`
}

type OrderListService struct {
	Status   int `form:"status" json:"status"`
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

func (service *CreateOrderService) Create(userID uint) dto.Response {
	code := e.SUCCESS

	var address model.Address
	if err := db.GetDB().Where("id = ? AND user_id = ?", service.AddressID, userID).First(&address).Error; err != nil {
		return dto.Response{
			Status: e.ErrorNotExistAddress,
			Msg:    e.GetMsg(e.ErrorNotExistAddress),
		}
	}

	var carts []model.Cart
	if err := db.GetDB().Where("user_id = ?", userID).Find(&carts).Error; err != nil || len(carts) == 0 {
		return dto.Response{
			Status: e.InvalidParams,
			Msg:    "购物车为空",
		}
	}

	var totalPrice float64
	type cartProduct struct {
		Cart    model.Cart
		Product model.Product
	}
	var items []cartProduct

	for _, cart := range carts {
		var product model.Product
		if err := db.GetDB().First(&product, cart.ProductID).Error; err != nil {
			continue
		}
		if product.Num < int(cart.Num) {
			return dto.Response{
				Status: e.ERROR,
				Msg:    fmt.Sprintf("商品 %s 库存不足", product.Name),
			}
		}
		price, _ := strconv.ParseFloat(product.Price, 64)
		totalPrice += price * float64(cart.Num)
		items = append(items, cartProduct{Cart: cart, Product: product})
	}

	if len(items) == 0 {
		return dto.Response{
			Status: e.ERROR,
			Msg:    "没有有效的商品",
		}
	}

	for _, item := range items {
		lockKey := cache.GetProductLockKey(item.Product.ID)
		locked, err := cache.GetInstance().Lock(lockKey, 5*time.Second)
		if err != nil || !locked {
			return dto.Response{
				Status: e.ERROR,
				Msg:    fmt.Sprintf("商品 %s 正在被抢购，请稍后重试", item.Product.Name),
			}
		}
		defer cache.GetInstance().Unlock(lockKey)

		result := db.GetDB().Model(&model.Product{}).
			Where("id = ? AND num >= ?", item.Product.ID, item.Cart.Num).
			Update("num", gorm.Expr("num - ?", item.Cart.Num))
		if result.RowsAffected == 0 {
			return dto.Response{
				Status: e.ERROR,
				Msg:    fmt.Sprintf("商品 %s 库存不足", item.Product.Name),
			}
		}
	}

	orderNo := generateOrderNo()
	order := model.Order{
		OrderNo:    orderNo,
		UserID:     userID,
		AddressID:  service.AddressID,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPending,
		Remark:     service.Remark,
	}

	if err := db.GetDB().Create(&order).Error; err != nil {
		pkg_logger.Logger.Error("create order error", "error", err)
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	var orderItems []model.OrderItem
	for _, item := range items {
		price, _ := strconv.ParseFloat(item.Product.Price, 64)
		orderItem := model.OrderItem{
			OrderID:     order.ID,
			ProductID:   item.Product.ID,
			ProductName: item.Product.Name,
			Num:         item.Cart.Num,
			Price:       price,
		}
		orderItems = append(orderItems, orderItem)
	}
	db.GetDB().Create(&orderItems)

	db.GetDB().Where("user_id = ?", userID).Delete(&model.Cart{})

	return dto.Response{
		Status: code,
		Data:   dto.BuildOrder(order, orderItems),
		Msg:    "下单成功",
	}
}

func (service *OrderListService) List(userID uint) dto.Response {
	var orders []model.Order
	var total int64
	code := e.SUCCESS

	if service.PageSize == 0 {
		service.PageSize = 10
	}
	if service.Page == 0 {
		service.Page = 1
	}

	query := db.GetDB().Model(&model.Order{}).Where("user_id = ?", userID)
	if service.Status > 0 {
		query = query.Where("status = ?", service.Status)
	}

	query.Count(&total)

	offset := (service.Page - 1) * service.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(service.PageSize).Find(&orders).Error; err != nil {
		pkg_logger.Logger.Error("list orders error", "error", err)
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return dto.BuildPagedResponse(dto.BuildOrders(orders), total, service.Page, service.PageSize)
}

func OrderDetail(orderID string, userID uint) dto.Response {
	var order model.Order
	code := e.SUCCESS

	if err := db.GetDB().Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		return dto.Response{
			Status: e.ErrorDatabase,
			Msg:    "订单不存在",
		}
	}

	var items []model.OrderItem
	db.GetDB().Where("order_id = ?", order.ID).Find(&items)

	return dto.Response{
		Status: code,
		Data:   dto.BuildOrder(order, items),
		Msg:    e.GetMsg(code),
	}
}

func OrderPay(orderID string, userID uint) dto.Response {
	var order model.Order
	if err := db.GetDB().Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		return dto.Response{
			Status: e.ErrorDatabase,
			Msg:    "订单不存在",
		}
	}

	if order.Status != model.OrderStatusPending {
		return dto.Response{
			Status: e.ERROR,
			Msg:    "订单状态不正确，无法支付",
		}
	}

	order.Status = model.OrderStatusPaid
	if err := db.GetDB().Save(&order).Error; err != nil {
		pkg_logger.Logger.Error("pay order error", "error", err)
		return dto.Response{
			Status: e.ErrorDatabase,
			Msg:    e.GetMsg(e.ErrorDatabase),
			Error:  err.Error(),
		}
	}

	return dto.Response{
		Status: e.SUCCESS,
		Msg:    "支付成功",
	}
}

func OrderCancel(orderID string, userID uint) dto.Response {
	var order model.Order
	if err := db.GetDB().Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		return dto.Response{
			Status: e.ErrorDatabase,
			Msg:    "订单不存在",
		}
	}

	if order.Status != model.OrderStatusPending {
		return dto.Response{
			Status: e.ERROR,
			Msg:    "订单状态不正确，无法取消",
		}
	}

	var items []model.OrderItem
	db.GetDB().Where("order_id = ?", order.ID).Find(&items)
	for _, item := range items {
		db.GetDB().Model(&model.Product{}).
			Where("id = ?", item.ProductID).
			Update("num", gorm.Expr("num + ?", item.Num))
	}

	order.Status = model.OrderStatusCancelled
	if err := db.GetDB().Save(&order).Error; err != nil {
		pkg_logger.Logger.Error("cancel order error", "error", err)
		return dto.Response{
			Status: e.ErrorDatabase,
			Msg:    e.GetMsg(e.ErrorDatabase),
			Error:  err.Error(),
		}
	}

	return dto.Response{
		Status: e.SUCCESS,
		Msg:    "订单已取消，库存已恢复",
	}
}

func generateOrderNo() string {
	return time.Now().Format("20060102150405") + uuid.New().String()[:8]
}
