package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
{
	"code": 1001, // 程序的错误码（和前端约定）
	"msg": xx, // 提示信息
	"data": {} // 数据
}
*/

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func ResponseError(c *gin.Context, code ResCode) {
	rd := &ResponseData{
		Code: code,
		Msg:  code.GetMsg(),
		Data: nil,
	}

	c.JSON(http.StatusOK, rd)
}

func ResponseErrorWithMsg(ctx *gin.Context, code ResCode, msg interface{}) {
	rd := &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	}

	ctx.JSON(http.StatusOK, rd)
}

func ResponseSuccess(ctx *gin.Context, data interface{}) {
	rd := &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.GetMsg(),
		Data: data,
	}

	ctx.JSON(http.StatusOK, rd)
}
