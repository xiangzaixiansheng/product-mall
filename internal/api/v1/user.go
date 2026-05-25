package v1

import (
	"net/http"
	"product-mall/internal/service"
	"product-mall/pkg/pkg_logger"

	"github.com/gin-gonic/gin"
)

// UserRegister 用户注册
// @Summary 用户注册
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body service.UserService true "注册信息"
// @Success 200 {object} dto.Response
// @Router /user/register [post]
func UserRegister(c *gin.Context) {
	var userRegisterService service.UserService
	if err := c.ShouldBind(&userRegisterService); err == nil {
		res := userRegisterService.Register()
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		pkg_logger.Logger.Error("error", "error", err)
	}
}

// UserLogin 用户登录
// @Summary 用户登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body service.UserService true "登录信息"
// @Success 200 {object} dto.Response
// @Router /user/login [post]
func UserLogin(c *gin.Context) {
	var userLoginService service.UserService
	if err := c.ShouldBind(&userLoginService); err == nil {
		res := userLoginService.Login()
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		pkg_logger.Logger.Error("error", "error", err)
	}
}

// UserUpdate 更新用户信息
// @Summary 更新用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body service.UserService true "用户信息"
// @Success 200 {object} dto.Response
// @Router /user [put]
func UserUpdate(c *gin.Context) {
	var userUpdateService service.UserService
	if err := c.ShouldBind(&userUpdateService); err == nil {
		res := userUpdateService.Update(getUserID(c))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		pkg_logger.Logger.Error("error", "error", err)
	}
}

// SendEmail 发送邮件
// @Summary 发送验证邮件
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body service.SendEmailService true "邮件信息"
// @Success 200 {object} dto.Response
// @Router /user/sending-email [post]
func SendEmail(c *gin.Context) {
	var sendEmailService service.SendEmailService
	if err := c.ShouldBind(&sendEmailService); err == nil {
		res := sendEmailService.SendEmail(getUserID(c))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		pkg_logger.Logger.Error("error", "error", err)
	}
}
