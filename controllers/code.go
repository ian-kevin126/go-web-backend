package controllers

/**
 * @Author ian-kevin
 * @Description //TODO 定义业务状态码
 * @Date 22:11 2022/2/10
 **/

type ResCode int64

const (
	CodeSuccess         ResCode = 1000
	CodeInvalidParams   ResCode = 1001
	CodeUserExist       ResCode = 1002
	CodeUserNotExist    ResCode = 1003
	CodeInvalidPassword ResCode = 1004
	CodeServerBusy      ResCode = 1005

	CodeInvalidToken      ResCode = 1006
	CodeInvalidAuthFormat ResCode = 1007
	CodeNotLogin          ResCode = 1008
)

var msgFlags = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParams:   "请求参数错误",
	CodeUserExist:       "用户名重复",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeInvalidToken:      "无效的Token",
	CodeInvalidAuthFormat: "认证格式有误",
	CodeNotLogin:          "未登录",
}

func (c ResCode) GetMsg() string {
	msg, ok := msgFlags[c]
	if ok {
		return msg
	}
	return msgFlags[CodeServerBusy]
}
