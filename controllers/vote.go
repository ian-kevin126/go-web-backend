package controllers

import (
	"fmt"
	"gin_demo/logic"
	"gin_demo/models"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

func PostVoteHandler(c *gin.Context) {
	// 1，参数校验
	vote := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(&vote); err != nil {
		errs, ok := err.(validator.ValidationErrors) // 类型断言
		if !ok {
			ResponseError(c, CodeInvalidParams)
			return
		}
		errData := removeTopStruct(errs.Translate(trans)) // 翻译并去除掉错误提示中的结构体标识
		ResponseErrorWithMsg(c, CodeInvalidParams, errData)
		return
	}

	// 获取当前请求用户的userid
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}

	fmt.Println("ppp", vote)
	// 具体投票的业务逻辑
	if err := logic.VoteForPost(userID, vote); err != nil {
		zap.L().Error("logic.VoteForPost(userID, p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}
