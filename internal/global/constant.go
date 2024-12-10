package global

const (
	Success             = 0
	ParamsNull          = 1
	ParamsError         = 2
	SignatureError      = 3
	RequestTimeOut      = 4
	InternalServerError = 5
	DeafultFailed       = 8 // 该错误通常要重新发起请求

	ValidateError = 1201

	TokenExpired  = 1001
	TokenNotFound = 1002
	TokenError    = 1003
	Unknown       = 666
)
