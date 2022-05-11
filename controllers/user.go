package controllers

import (
	"errors"
	"fmt"
	"gin_demo/dao/mysql"
	"gin_demo/logic"
	"gin_demo/models"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func SignUpHandler(c *gin.Context) {
	// 1，获取参数和参数的校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回错误相应
		zap.L().Error("signup with invalid param", zap.Error(err))
		// 判断err是否是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParams)
			return
		}
		// 如果是，则将报错翻译成中文再返回
		ResponseErrorWithMsg(c, CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(errs.Translate(trans)),
		//})
		return
	}
	// 手动对请求参数进行详细的业务规则校验
	/*
		if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.RePassword != p.Password {
			zap.L().Error("signup with invalid param")
			c.JSON(http.StatusOK, gin.H{
				"msg": "请求参数有误",
			})
			return
		}
	*/

	// 使用第三方库validator来校验参数，通过在models.ParamSignUp中使用binding tag来实现上述的字段required。

	fmt.Println(p)

	// 2，业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.Signup failed", zap.Error(err))
		// 如果用户已存在
		if errors.Is(err, mysql.ErrorUserExit) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3，返回响应
	ResponseSuccess(c, nil)
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "注册成功",
	//})
}

func LoginHandler(c *gin.Context) {
	// 1，获取请求参数
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParams)
			//c.JSON(http.StatusOK, gin.H{
			//	"msg": err.Error(),
			//})
			return
		}

		ResponseErrorWithMsg(c, CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(errs.Translate(trans)),
		//})
		return
	}

	// 2，业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("Logic.Login failed", zap.String("username", p.Username), zap.Error(err))

		if errors.Is(err, mysql.ErrorUserNotExit) {
			ResponseError(c, CodeUserNotExist)
			return
		}

		ResponseError(c, CodeInvalidPassword)
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "用户或者密码错误",
		//})
		return
	}

	// 3，返回响应
	ResponseSuccess(c, gin.H{
		"userId":   fmt.Sprintf("%d", user.UserID), // 如果id值大于js中能表示的最大值（2的53次方减一，而go中的int64最大可以表示到2的63次方），就会出现失真
		"userName": user.UserName,
		"token":    user.Token,
	})
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "登录成功",
	//})
}
