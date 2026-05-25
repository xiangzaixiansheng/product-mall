package service

import (
	"mime/multipart"
	"product-mall/internal/dto"
	"product-mall/internal/model"
	"product-mall/pkg/db"
	"product-mall/pkg/e"
	"product-mall/pkg/pkg_logger"
)

/*
创建商品信息
**/

type ProductService struct {
	ID            uint   `form:"id" json:"id"`
	Name          string `form:"name" json:"name"`
	CategoryID    int    `form:"category_id" json:"category_id"`
	Title         string `form:"title" json:"title" binding:"required,min=2,max=100"`
	Info          string `form:"info" json:"info" binding:"max=1000"`
	ImgPath       string `form:"img_path" json:"img_path"`
	Price         string `form:"price" json:"price"`
	DiscountPrice string `form:"discount_price" json:"discount_price"`
	OnSale        bool   `form:"on_sale" json:"on_sale"`
	Num           int    `form:"num" json:"num"`
	PageNum       int    `form:"pageNum"`
	PageSize      int    `form:"pageSize"`
}

type ListProductImgService struct {
}

//创建商品
func (service *ProductService) Create(id uint, files []*multipart.FileHeader) dto.Response {
	code := e.SUCCESS
	//获取用户信息
	var user model.User
	if err := db.GetDB().Model(&model.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		code = e.ErrorExistUser
		return dto.Response{
			Status: code,
			Data:   e.GetMsg(code),
		}
	}
	tmp, _ := files[0].Open()
	status, info := Upload2QiNiu(tmp, files[0].Size)
	if status != 200 {
		return dto.Response{
			Status: status,
			Data:   e.GetMsg(status),
			Error:  info,
		}
	}
	//存储product
	product := model.Product{
		Name:             service.Name,
		CategoryID:       uint(service.CategoryID),
		Title:            service.Title,
		Info:             service.Info,
		ImgPath:          info,
		Price:            service.Price,
		DiscountPrice:    service.DiscountPrice,
		Num:              service.Num,
		OnSale:           true,
		CreateUserID:     int(id),
		CreateUserName:   user.UserName,
		CreateUserAvatar: user.Avatar,
	}
	err := db.GetDB().Create(&product).Error
	if err != nil {
		pkg_logger.LogrusObj.Error("error", "error", err)
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//解析文件信息--其他文件存数据库-商品id和文件地址保存在一起
	for _, file := range files {
		tmp, _ := file.Open()
		status, info := Upload2QiNiu(tmp, file.Size)
		if status != 200 {
			return dto.Response{
				Status: status,
				Data:   e.GetMsg(status),
				Error:  info,
			}
		}
		// 每一项图片都存在数据库里面
		productImg := model.ProductImg{
			ProductID: product.ID,
			ImgPath:   info,
		}
		err = db.GetDB().Create(&productImg).Error
		if err != nil {
			code = e.ERROR
			return dto.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
	}
	return dto.Response{
		Status: code,
		Data:   dto.BuildProduct(product),
		Msg:    e.GetMsg(code),
	}

}

//List接口
func (service *ProductService) List() dto.Response {
	var products []model.Product
	var total int64
	code := e.SUCCESS
	if service.PageSize == 0 {
		service.PageSize = 20
	}
	//如果传入的商品的CategoryID为0的话
	if service.CategoryID == 0 {
		if err := db.GetDB().Model(model.Product{}).Count(&total).Error; err != nil {
			pkg_logger.LogrusObj.Error("error", "error", err)
			code = e.ErrorDatabase
			return dto.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		if err := db.GetDB().Offset((service.PageNum - 1) * service.PageSize).
			Limit(service.PageSize).Find(&products).
			Error; err != nil {
			pkg_logger.LogrusObj.Error("error", "error", err)
			code = e.ErrorDatabase
			return dto.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

	} else {
		//Preload 预处理用来处理一对多关系的
		if err := db.GetDB().Model(model.Product{}).Preload("Category").
			Where("category_id = ?", service.CategoryID).
			Count(&total).Error; err != nil {
			pkg_logger.LogrusObj.Error("error", "error", err)
			code = e.ErrorDatabase
			return dto.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		if err := db.GetDB().Model(model.Product{}).Preload("Category").
			Where("category_id=?", service.CategoryID).
			Offset((service.PageNum - 1) * service.PageSize).
			Limit(service.PageSize).
			Find(&products).Error; err != nil {
			pkg_logger.LogrusObj.Error("error", "error", err)
			code = e.ErrorDatabase
			return dto.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

	}

	return dto.BuildListResponse(dto.BuildProducts(products), uint(total))
}

//删除商品
func (service *ProductService) Delete(id string) dto.Response {
	code := e.SUCCESS
	var product model.Product
	//判断商品是否存在
	if err := db.GetDB().First(&product, id).Error; err != nil {
		pkg_logger.LogrusObj.Error("error", "error", err)
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//存在则删除商品
	if err := db.GetDB().Delete(&product).Error; err != nil {
		pkg_logger.LogrusObj.Error("error", "error", err)
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

//更新商品
func (service *ProductService) Update(id string) dto.Response {
	var product model.Product
	db.GetDB().Model(&model.Product{}).First(&product, id)
	product.Name = service.Name
	product.CategoryID = uint(service.CategoryID)
	product.Title = service.Title
	product.Info = service.Info
	product.ImgPath = service.ImgPath
	product.Price = service.Price
	product.DiscountPrice = service.DiscountPrice
	product.OnSale = service.OnSale
	code := e.SUCCESS

	if err := db.GetDB().Save(&product).Error; err != nil {
		pkg_logger.LogrusObj.Error("error", "error", err)
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

//搜索商品
func (service *ProductService) Search() dto.Response {
	var products []model.Product
	code := e.SUCCESS
	if service.PageSize == 0 {
		service.PageSize = 15
	}
	err := db.GetDB().Where("name LIKE ? OR info LIKE ?", "%"+service.Info+"%", "%"+service.Info+"%").
		Offset((service.PageNum - 1) * service.PageSize).
		Limit(service.PageSize).Find(&products).Error
	if err != nil {
		pkg_logger.LogrusObj.Error("error", "error", err)
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return dto.BuildListResponse(dto.BuildProducts(products), uint(len(products)))
}

//获取商品列表图片
func (service *ListProductImgService) List(id string) dto.Response {
	var productImgList []model.ProductImg
	db.GetDB().Model(model.ProductImg{}).Where("product_id=?", id).Find(&productImgList)
	return dto.BuildListResponse(dto.BuildProductImgs(productImgList), uint(len(productImgList)))
}
