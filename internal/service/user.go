package service

import (
	"product-mall/conf"
	"product-mall/internal/dto"
	"product-mall/internal/model"
	util "product-mall/internal/tools"
	"product-mall/pkg/db"
	"product-mall/pkg/e"
	"product-mall/pkg/pkg_logger"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
	"gopkg.in/gomail.v2"
)

//UserRegisterService 管理用户注册服务
/**
用户注册
用户登录
修改用户名字
绑定邮箱
验证邮箱
*/
type UserService struct {
	NickName string `form:"nick_name" json:"nick_name"`
	UserName string `form:"user_name" json:"user_name"`
	Password string `form:"password" json:"password"`
	Key      string `form:"key" json:"key"` // 前端进行判断
}

type SendEmailService struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
	//OpertionType 1:绑定邮箱 2：解绑邮箱 3：改密码
	OperationType uint `form:"operation_type" json:"operation_type"`
}

type ValidEmailService struct {
}

//注册
func (service UserService) Register() dto.Response {
	var user model.User
	var count int64
	code := e.SUCCESS
	if service.Key == "" || len(service.Key) != 16 {
		code = e.ERROR
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "密钥长度不足",
		}
	}
	//密钥，支付密码
	util.Encrypt.SetKey(service.Key)
	db.GetDB().Model(&model.User{}).Where("user_name=?", service.UserName).Count(&count)
	if count == 1 {
		code = e.ErrorExistUser
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	user = model.User{
		Nickname: service.NickName,
		UserName: service.UserName,
		Status:   model.Active,
		Money:    util.Encrypt.AesEncoding("10000"),
	}
	if err := user.SetPassword(service.Password); err != nil {
		pkg_logger.Logger.Error("error", "error", err)
		code = e.ErrorFailEncryption
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	user.Avatar = "http://q1.qlogo.cn/g?b=qq&nk=294350394&s=640"
	//创建用户
	if err := db.GetDB().Create(&user).Error; err != nil {
		pkg_logger.Logger.Error("error", "error", err)
		code = e.ErrorDatabase
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

//用户登陆函数
func (service UserService) Login() dto.Response {
	var user model.User
	code := e.SUCCESS
	if err := db.GetDB().Where("user_name=?", service.UserName).First(&user).Error; err != nil {
		pkg_logger.Logger.Error("login error", "error", err)

		//如果查询不到，返回相应的错误
		if errors.Is(err, gorm.ErrRecordNotFound) {

			code = e.ErrorNotExistUser
			return dto.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	if user.CheckPassword(service.Password) == false {
		code = e.ErrorNotCompare
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	token, err := util.GenerateToken(user.ID, service.UserName, 0)
	if err != nil {
		pkg_logger.Logger.Error("error", "error", err)
		code = e.ErrorAuthToken
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return dto.Response{
		Status: code,
		Data:   dto.TokenData{User: dto.BuildUser(user), Token: token},
		Msg:    e.GetMsg(code),
	}
}

func (service UserService) Update(id uint) dto.Response {
	var user model.User
	code := e.SUCCESS
	//https://gorm.io/zh_CN/docs/query.html
	err := db.GetDB().Where("id", id).First(&user).Error
	if err != nil {
		pkg_logger.Logger.Error("error", "error", err)
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	if service.NickName != "" {
		user.Nickname = service.NickName
	}

	err = db.GetDB().Save(&user).Error
	if err != nil {
		pkg_logger.Logger.Error("error", "error", err)
		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return dto.Response{
		Status: code,
		Data:   dto.BuildUser(user),
		Msg:    e.GetMsg(code),
	}

}

//检查email中的token
func (service UserService) Valid_token(token string) dto.Response {
	var userID uint
	var email string
	var password string
	//操作类型
	var operationType uint
	//用户信息
	var user model.User

	code := e.SUCCESS

	if token == "" {
		code = e.InvalidParams
	} else {
		//解析链接中的token
		claims, err := util.ParseEmailToken(token)
		if err != nil {
			pkg_logger.Logger.Error("error", "error", err)
			code = e.ErrorAuthCheckTokenFail
		} else if claims.ExpiresAt.Before(time.Now()) {
			code = e.ErrorAuthCheckTokenTimeout
		} else {
			userID = claims.UserID
			email = claims.Email
			password = claims.Password
			operationType = claims.OperationType
		}
	}

	if code != e.SUCCESS {
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 绑定邮箱
	if operationType == 1 {
		if err := db.GetDB().Table("user").Where("id=?", userID).Update("email", email).Error; err != nil {
			pkg_logger.Logger.Error("valid_token bind-email error", "error", err)
			code = e.ErrorDatabase
			return dto.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
	} else if operationType == 2 {
		//解除绑定邮箱
		if err := db.GetDB().Table("user").Where("id=?", userID).Update("email", "").Error; err != nil {
			pkg_logger.Logger.Error("valid_token unbind-email error", "error", err)
			code = e.ErrorDatabase
			return dto.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
	} else if operationType == 3 {
		//获取用户信息
		if err := db.GetDB().First(&user, userID).Error; err != nil {
			pkg_logger.Logger.Error("error", "error", err)
			code = e.ErrorDatabase
			return dto.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
		// 对密码进行加密
		if err := user.SetPassword(password); err != nil {
			pkg_logger.Logger.Error("error", "error", err)
			code = e.ErrorFailEncryption
			return dto.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
		//更新数据
		if err := db.GetDB().Save(&user).Error; err != nil {
			pkg_logger.Logger.Error("error", "error", err)
			code = e.ErrorDatabase
			return dto.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
		code = e.UpdatePasswordSuccess
		//返回修改密码成功信息
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	} else {
		//没有匹配的方法
		code = e.InvalidParams
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	//返回用户信息
	return dto.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   dto.BuildUser(user),
	}

}

//发送邮件的接口, 方法传入id为用户信息的id
func (service SendEmailService) SendEmail(id uint) dto.Response {
	code := e.SUCCESS
	var noticeMsg model.Notice
	var emailAddress string

	token, err := util.GenerateEmailToken(id, service.OperationType, service.Password, service.Email)
	if err != nil {
		pkg_logger.Logger.Error("send_email generate_token error", "error", err)
		code = e.ErrorAuthToken
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//获取对应的邮件提醒的数据
	if err := db.GetDB().First(&noticeMsg, service.OperationType).Error; err != nil {
		pkg_logger.Logger.Error("send_email get_email_data error", "error", err)

		code = e.ErrorDatabase
		return dto.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	//替换邮件内容,准备发送邮件
	//服务验证email的地址 + token
	emailAddress = conf.ValidEmail + token
	mailTitle := noticeMsg.Text
	mailText := strings.Replace(mailTitle, "Email", emailAddress, -1)

	//发送邮件
	m := gomail.NewMessage()
	m.SetHeader("From", conf.SmtpEmail)
	m.SetHeader("To", service.Email)
	m.SetHeader("Subject", "Gin-tool") //邮件标题
	m.SetBody("text/html", mailText)   //邮件内容
	//m.Attach("E:\\IMGP0814.JPG")   //邮件附件
	d := gomail.NewDialer(conf.SmtpHost, 465, conf.SmtpEmail, conf.SmtpPass)

	if err := d.DialAndSend(m); err != nil {
		pkg_logger.Logger.Error("send_email error", "error", err)
		code = e.ErrorSendEmail
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
