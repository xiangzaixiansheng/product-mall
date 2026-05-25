package service

import (
	"product-mall/internal/dto"
	"product-mall/internal/model"
	"product-mall/pkg/db"
	"product-mall/pkg/e"
	"product-mall/pkg/pkg_logger"
)

type CategoryService struct {
	Name     string `form:"name" json:"name" binding:"required,min=1,max=100"`
	ParentID uint   `form:"parent_id" json:"parent_id"`
	Sort     int    `form:"sort" json:"sort"`
}

func (service *CategoryService) Create() dto.Response {
	code := e.SUCCESS
	var level uint = 1

	if service.ParentID > 0 {
		var parent model.Category
		if err := db.GetDB().First(&parent, service.ParentID).Error; err != nil {
			code = e.ErrorDatabase
			return dto.Response{
				Status: code,
				Msg:    "父分类不存在",
			}
		}
		level = parent.Level + 1
	}

	category := model.Category{
		Name:     service.Name,
		ParentID: service.ParentID,
		Level:    level,
		Sort:     service.Sort,
	}

	if err := db.GetDB().Create(&category).Error; err != nil {
		pkg_logger.Logger.Error("create category error", "error", err)
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return dto.Response{
		Status: code,
		Data:   dto.BuildCategory(category),
		Msg:    e.GetMsg(code),
	}
}

func (service *CategoryService) List() dto.Response {
	var categories []model.Category
	code := e.SUCCESS

	if err := db.GetDB().Order("sort ASC, id ASC").Find(&categories).Error; err != nil {
		pkg_logger.Logger.Error("list categories error", "error", err)
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return dto.Response{
		Status: code,
		Data:   dto.BuildCategoryTree(categories),
		Msg:    e.GetMsg(code),
	}
}

func (service *CategoryService) Update(id string) dto.Response {
	var category model.Category
	code := e.SUCCESS

	if err := db.GetDB().First(&category, id).Error; err != nil {
		pkg_logger.Logger.Error("find category error", "error", err)
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    "分类不存在",
		}
	}

	category.Name = service.Name
	category.Sort = service.Sort

	if err := db.GetDB().Save(&category).Error; err != nil {
		pkg_logger.Logger.Error("update category error", "error", err)
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return dto.Response{
		Status: code,
		Data:   dto.BuildCategory(category),
		Msg:    e.GetMsg(code),
	}
}

func (service *CategoryService) Delete(id string) dto.Response {
	code := e.SUCCESS

	var count int64
	db.GetDB().Model(&model.Category{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		return dto.Response{
			Status: e.ERROR,
			Msg:    "该分类下有子分类，无法删除",
		}
	}

	if err := db.GetDB().Delete(&model.Category{}, id).Error; err != nil {
		pkg_logger.Logger.Error("delete category error", "error", err)
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return dto.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
