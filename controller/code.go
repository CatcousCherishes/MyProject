package controller

type ResCode int64

// 常量
const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodeNeedLogin
	CodeInvalidToken

	CodeNeesAuth
	CodeInvalidAuth
	CodeNeesLogin
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户名已存在",
	CodeUserNotExist:    "用户名不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",
	CodeNeedLogin:       "需要登录",
	CodeInvalidToken:    "无效的token",

	CodeNeesAuth:    "需要Auth",
	CodeInvalidAuth: "无效得Token",
	CodeNeesLogin:   "需要登录",
}

// Msg() code取出msg状态

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c] //根据c 取出msg
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
