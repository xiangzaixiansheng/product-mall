package service

import (
	"context"
	"product-mall/internal/dto"
	"product-mall/internal/model"
	"product-mall/internal/repo/mysql"
	"product-mall/pkg/e"
	"strconv"

	"product-mall/pkg/pkg_logger"
)

//绑定json数据
type AddressService struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func (service AddressService) Create(ctx context.Context, id uint) dto.Response {
	//插入数据
	code := e.SUCCESS
	repo := mysql.NewAddressRepo(ctx)

	address := model.Address{
		UserID:  id,
		Name:    service.Name,
		Phone:   service.Phone,
		Address: service.Address,
	}
	err := repo.Create(&address)
	if err != nil {
		pkg_logger.Logger.Error("db error", "error", err)
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//返回数据库中这个用户最新的地址信息
	var addresses []model.Address
	addresses, err = repo.GetAddressByUid(id)

	if err != nil {
		code = e.ErrorDatabase
		pkg_logger.Logger.Error("db error", "error", err)
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return dto.Response{
		Status: code,
		Data:   dto.BuildAddresses(addresses),
		Msg:    e.GetMsg(code),
	}

}

func (service AddressService) List(ctx context.Context, id uint) dto.Response {
	code := e.SUCCESS
	repo := mysql.NewAddressRepo(ctx)
	var addresses []model.Address
	addresses, err := repo.GetAddressByUid(id)
	if err != nil {
		pkg_logger.Logger.Error("db error", "error", err)
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return dto.Response{
		Status: code,
		Data:   dto.BuildAddresses(addresses),
		Msg:    e.GetMsg(code),
	}

}

func (service AddressService) Delete(ctx context.Context, id string) dto.Response {
	//可以先找出来对应id的数据
	code := e.SUCCESS
	repo := mysql.NewAddressRepo(ctx)

	var address model.Address
	address, err := repo.GetAddressById(id)
	if err != nil {
		pkg_logger.Logger.Error("db error", "error", err)
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	err = repo.DeleteAddress(address)
	if err != nil {
		code = e.ErrorDatabase
		pkg_logger.Logger.Error("error", "error", err)
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	return dto.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}

}
func (service AddressService) Update(ctx context.Context, uid uint, aid string) dto.Response {
	code := e.SUCCESS
	repo := mysql.NewAddressRepo(ctx)

	address := model.Address{
		Address: service.Address,
		Phone:   service.Phone,
		Name:    service.Name,
		UserID:  uid,
	}

	Id, _ := strconv.Atoi(aid)
	address.ID = uint(Id)

	//字段全部更新
	err := repo.Updates(&address)
	if err != nil {
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	var addresses []model.Address
	addresses, err = repo.GetAddressByUid(uid)
	if err != nil {
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	return dto.Response{
		Status: code,
		Data:   dto.BuildAddresses(addresses),
	}

}

func (service *AddressService) Show(ctx context.Context, id string) dto.Response {
	var addresses []model.Address
	repo := mysql.NewAddressRepo(ctx)

	code := e.SUCCESS
	addresses, err := repo.GetAddressByUid(id)

	if err != nil {
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return dto.Response{
		Status: code,
		Data:   dto.BuildAddresses(addresses),
		Msg:    e.GetMsg(code),
	}

}
