package v1

import (
	"product-mall/internal/service"
	"product-mall/pkg/pkg_logger"

	"github.com/gin-gonic/gin"
)

func CreateAddress(c *gin.Context) {
	svc := service.AddressService{}
	if err := c.BindJSON(&svc); err == nil {
		res := svc.Create(c, getUserID(c))
		c.JSON(200, res)
	} else {
		pkg_logger.Logger.Error("error", "error", err)
		c.JSON(400, ErrorResponse(err))
	}
}

func ShowAddresses(c *gin.Context) {
	svc := service.AddressService{}
	res := svc.Show(c, c.Param("id"))
	c.JSON(200, res)
}

func UpdateAddress(c *gin.Context) {
	svc := service.AddressService{}
	if err := c.BindJSON(&svc); err == nil {
		res := svc.Update(c, getUserID(c), c.Param("id"))
		c.JSON(200, res)
	} else {
		pkg_logger.Logger.Error("error", "error", err)
		c.JSON(400, ErrorResponse(err))
	}
}

func DeleteAddress(c *gin.Context) {
	svc := service.AddressService{}
	res := svc.Delete(c, c.Param("id"))
	c.JSON(200, res)
}
