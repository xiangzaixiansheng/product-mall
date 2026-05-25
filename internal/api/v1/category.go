package v1

import (
	"net/http"
	"product-mall/internal/service"
	"product-mall/pkg/pkg_logger"

	"github.com/gin-gonic/gin"
)

// ListCategories 获取分类树
// @Summary 获取商品分类树
// @Tags 分类
// @Produce json
// @Success 200 {object} dto.Response
// @Router /categories [get]
func ListCategories(c *gin.Context) {
	svc := service.CategoryService{}
	res := svc.List()
	c.JSON(http.StatusOK, res)
}

// CreateCategory 创建分类
// @Summary 创建商品分类
// @Tags 分类
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body service.CategoryService true "分类信息"
// @Success 200 {object} dto.Response
// @Router /categories [post]
func CreateCategory(c *gin.Context) {
	var svc service.CategoryService
	if err := c.ShouldBindJSON(&svc); err == nil {
		res := svc.Create()
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		pkg_logger.Logger.Error("error", "error", err)
	}
}

// UpdateCategory 更新分类
// @Summary 更新商品分类
// @Tags 分类
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "分类ID"
// @Param body body service.CategoryService true "分类信息"
// @Success 200 {object} dto.Response
// @Router /categories/{id} [put]
func UpdateCategory(c *gin.Context) {
	var svc service.CategoryService
	if err := c.ShouldBindJSON(&svc); err == nil {
		res := svc.Update(c.Param("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		pkg_logger.Logger.Error("error", "error", err)
	}
}

// DeleteCategory 删除分类
// @Summary 删除商品分类
// @Tags 分类
// @Produce json
// @Security BearerAuth
// @Param id path string true "分类ID"
// @Success 200 {object} dto.Response
// @Router /categories/{id} [delete]
func DeleteCategory(c *gin.Context) {
	svc := service.CategoryService{}
	res := svc.Delete(c.Param("id"))
	c.JSON(http.StatusOK, res)
}
