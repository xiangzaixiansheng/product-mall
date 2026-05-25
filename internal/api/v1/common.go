package v1

import (
	"encoding/json"
	"fmt"
	"product-mall/conf"
	"product-mall/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func getUserID(c *gin.Context) uint {
	id, _ := c.Get("user_id")
	uid, _ := id.(uint)
	return uid
}

func ErrorResponse(err error) dto.Response {
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			field := conf.T(fmt.Sprintf("Field.%s", e.Field()))
			tag := conf.T(fmt.Sprintf("Tag.Valid.%s", e.Tag()))
			return dto.Response{
				Status: 400,
				Msg:    fmt.Sprintf("%s%s", field, tag),
				Error:  fmt.Sprint(err),
			}
		}
	}
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return dto.Response{
			Status: 400,
			Msg:    "JSON类型不匹配",
			Error:  fmt.Sprint(err),
		}
	}

	return dto.Response{
		Status: 400,
		Msg:    "参数错误",
		Error:  fmt.Sprint(err),
	}
}
