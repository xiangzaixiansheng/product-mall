package v1

import (
	"net/http"
	"product-mall/internal/service"
	"product-mall/pkg/pkg_logger"

	"github.com/gin-gonic/gin"
)

func CreateCart(c *gin.Context) {
	createCartService := service.CartService{}
	if err := c.ShouldBind(&createCartService); err == nil {
		res := createCartService.Create(c.Param("id"), getUserID(c))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		pkg_logger.Logger.Error("error", "error", err)
	}
}

func ShowList(c *gin.Context) {
	showCartsService := service.CartService{}
	res := showCartsService.List(c.Param("id"))
	c.JSON(http.StatusOK, res)
}

func UpdateCart(c *gin.Context) {
	updateCartService := service.CartService{}
	if err := c.ShouldBind(&updateCartService); err == nil {
		res := updateCartService.Update(c.Param("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		pkg_logger.Logger.Error("error", "error", err)
	}
}

func DeleteCart(c *gin.Context) {
	deleteCartService := service.CartService{}
	if err := c.ShouldBind(&deleteCartService); err == nil {
		res := deleteCartService.Delete(c.Param("id"), getUserID(c))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		pkg_logger.Logger.Error("error", "error", err)
	}
}
