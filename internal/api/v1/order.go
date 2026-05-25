package v1

import (
	"net/http"
	"product-mall/internal/service"
	"product-mall/pkg/pkg_logger"

	"github.com/gin-gonic/gin"
)

// CreateOrder 创建订单
// @Summary 从购物车创建订单
// @Tags 订单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body service.CreateOrderService true "订单信息"
// @Success 200 {object} dto.Response
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	var svc service.CreateOrderService
	if err := c.ShouldBindJSON(&svc); err == nil {
		res := svc.Create(getUserID(c))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		pkg_logger.Logger.Error("error", "error", err)
	}
}

// ListOrders 订单列表
// @Summary 获取我的订单列表
// @Tags 订单
// @Produce json
// @Security BearerAuth
// @Param status query int false "订单状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} dto.Response
// @Router /orders [get]
func ListOrders(c *gin.Context) {
	var svc service.OrderListService
	if err := c.ShouldBindQuery(&svc); err == nil {
		res := svc.List(getUserID(c))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		pkg_logger.Logger.Error("error", "error", err)
	}
}

// GetOrder 订单详情
// @Summary 获取订单详情
// @Tags 订单
// @Produce json
// @Security BearerAuth
// @Param id path string true "订单ID"
// @Success 200 {object} dto.Response
// @Router /orders/{id} [get]
func GetOrder(c *gin.Context) {
	res := service.OrderDetail(c.Param("id"), getUserID(c))
	c.JSON(http.StatusOK, res)
}

// PayOrder 模拟支付
// @Summary 模拟支付订单
// @Tags 订单
// @Produce json
// @Security BearerAuth
// @Param id path string true "订单ID"
// @Success 200 {object} dto.Response
// @Router /orders/{id}/pay [put]
func PayOrder(c *gin.Context) {
	res := service.OrderPay(c.Param("id"), getUserID(c))
	c.JSON(http.StatusOK, res)
}

// CancelOrder 取消订单
// @Summary 取消订单并恢复库存
// @Tags 订单
// @Produce json
// @Security BearerAuth
// @Param id path string true "订单ID"
// @Success 200 {object} dto.Response
// @Router /orders/{id}/cancel [put]
func CancelOrder(c *gin.Context) {
	res := service.OrderCancel(c.Param("id"), getUserID(c))
	c.JSON(http.StatusOK, res)
}
