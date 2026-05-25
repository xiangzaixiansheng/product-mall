package v1

import (
	"net/http"
	"product-mall/internal/service"
	util "product-mall/internal/tools"
	"product-mall/pkg/pkg_logger"

	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	//相当于创建了一个UserRegisterService对象，调用这个对象中的Register方法。
	var userRegisterService service.UserService
	if err := c.ShouldBind(&userRegisterService); err == nil {
		res := userRegisterService.Register()
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		pkg_logger.LogrusObj.Error("error", "error", err)
	}
}

//UserLogin 用户登陆接口
func UserLogin(c *gin.Context) {
	var userLoginService service.UserService
	if err := c.ShouldBind(&userLoginService); err == nil {
		res := userLoginService.Login()
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		pkg_logger.LogrusObj.Error("error", "error", err)
	}
}

//更新用户信息
func UserUpdate(c *gin.Context) {
	var userUpdateService service.UserService
	claims, _ := util.ParseToken(c.GetHeader("Cookie"))
	if err := c.ShouldBind(&userUpdateService); err == nil {
		res := userUpdateService.Update(claims.ID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		pkg_logger.LogrusObj.Error("error", "error", err)
	}
}

//发送邮件
func SendEmail(c *gin.Context) {
	var sendEmailService service.SendEmailService
	//检查cookie里面的信息
	claims, _ := util.ParseToken(c.GetHeader("Cookie"))
	if err := c.ShouldBind(&sendEmailService); err == nil {
		res := sendEmailService.SendEmail(claims.ID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		pkg_logger.LogrusObj.Error("error", "error", err)
	}
}
