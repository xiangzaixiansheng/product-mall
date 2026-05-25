package e

const (
	SUCCESS               = 200
	UpdatePasswordSuccess = 201
	NotExistInentifier    = 202
	ERROR                 = 500
	InvalidParams         = 400

	//成员错误
	ErrorExistNick          = 10001
	ErrorExistUser          = 10002
	ErrorNotExistUser       = 10003
	ErrorNotCompare         = 10004
	ErrorNotComparePassword = 10005
	ErrorFailEncryption     = 10006
	ErrorNotExistProduct    = 10007
	ErrorNotExistAddress    = 10008
	ErrorExistFavorite      = 10009

	ErrorAuthCheckTokenFail        = 30001 //token 错误
	ErrorAuthCheckTokenTimeout     = 30002 //token 过期
	ErrorAuthToken                 = 30003
	ErrorAuth                      = 30004
	ErrorAuthInsufficientAuthority = 30005
	ErrorSendEmail                 = 30007

	//数据库错误
	ErrorDatabase = 40001

	//订单错误
	ErrorOrderNotFound    = 40010
	ErrorOrderStatus      = 40011
	ErrorOutOfStock       = 40012
	ErrorCartEmpty        = 40013

	//对象存储错误
	ErrorOss        = 50001
	ErrorUploadFile = 50002
)
