package v1

import (
	"log/slog"
	"net/http"
	"product-mall/internal/service"
	"product-mall/pkg/pkg_logger"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	var productService service.ProductService
	form, _ := c.MultipartForm()
	slog.Debug("multipart form received", "files", len(form.File["file"]))
	files := form.File["file"]
	if err := c.ShouldBind(&productService); err == nil {
		res := productService.Create(getUserID(c), files)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		pkg_logger.Logger.Error("error", "error", err)
	}
}

// SearchProducts 商品搜索
// @Summary 搜索商品列表
// @Tags 商品
// @Produce json
// @Param keyword query string false "搜索关键词"
// @Param category_id query int false "分类ID"
// @Param min_price query string false "最低价格"
// @Param max_price query string false "最高价格"
// @Param sort query string false "排序方式: price_asc/price_desc/newest/sales"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} dto.Response
// @Router /products [get]
func SearchProducts(c *gin.Context) {
	var svc service.ProductSearchService
	if err := c.ShouldBindQuery(&svc); err == nil {
		res := svc.Search()
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		pkg_logger.Logger.Error("error", "error", err)
	}
}
