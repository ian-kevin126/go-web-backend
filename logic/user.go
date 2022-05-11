package logic

import (
	"fmt"
	"gin_demo/dao/mysql"
	"gin_demo/models"
	"gin_demo/pkg/jwt"
	"gin_demo/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	// 1，判断用户是否存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		// 数据库查询出错
		return err
	}

	// 2，生成UID
	userID := snowflake.GenID()
	// 构造一个User示例
	user := &models.User{
		UserID:   userID,
		UserName: p.Username,
		Password: p.Password,
	}

	fmt.Println("user", user)

	// 3，保存用户到数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		UserName: p.Username,
		Password: p.Password,
	}

	// 传递的是指针，就能拿到票user.UserID
	if err := mysql.Login(user); err != nil {
		return nil, err
	}

	// 生成JWT
	token, err := jwt.GenToken(user.UserID, user.UserName)
	if err != nil {
		return
	}
	user.Token = token

	return
}
