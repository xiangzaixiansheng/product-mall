package v1

import (
	"product-mall/internal/service"
	util "product-mall/internal/tools"
	"product-mall/pkg/pkg_logger"

	"github.com/gin-gonic/gin"
)

//新增收货地址
func CreateAddress(c *gin.Context) {
	service := service.AddressService{}
	claim, _ := util.ParseToken(c.GetHeader("Cookie"))

	if err := c.BindJSON(&service); err == nil {
		res := service.Create(c, claim.ID)
		c.JSON(200, res)
	} else {
		pkg_logger.LogrusObj.Error("error", "error", err)
		c.JSON(400, ErrorResponse(err))
	}
}

//展示收货地址
func ShowAddresses(c *gin.Context) {
	service := service.AddressService{}
	res := service.Show(c, c.Param("id"))
	c.JSON(200, res)
}

//修改收货地址
func UpdateAddress(c *gin.Context) {
	service := service.AddressService{}
	claim, _ := util.ParseToken(c.GetHeader("Cookie"))
	if err := c.BindJSON(&service); err == nil {
		res := service.Update(c, claim.ID, c.Param("id"))
		c.JSON(200, res)
	} else {
		pkg_logger.LogrusObj.Error("error", "error", err)
		c.JSON(400, ErrorResponse(err))
	}
}

//删除收获地址
func DeleteAddress(c *gin.Context) {
	service := service.AddressService{}
	res := service.Delete(c, c.Param("id"))
	c.JSON(200, res)
}
