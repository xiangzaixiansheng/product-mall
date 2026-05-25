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
