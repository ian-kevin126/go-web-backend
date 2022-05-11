package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"gin_demo/models"
)

const secret = "dkahdkash_dsadas_JJSJ_#@J@J@H_sdasd"

// 把每一步数据库操作封装成函数
// 等待logic层根据业务需求调用

func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`

	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		fmt.Printf("db get user, err: %v\n", err)
		return err
	}

	if count > 0 {
		return ErrorUserExit
	}

	return
}

// InsertUser 向数据库中插入一条新的用户数据
func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)

	// 执行SQL语句入库
	sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.UserName, user.Password)
	return
}

// 用户密码加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	// 加盐字符串
	h.Write([]byte(secret))
	// 转化成16进制的字符串返回
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *models.User) (err error) {
	oPassword := user.Password // 用户登录的密码
	sqlStr := `select user_id, username, password from user where username=?`

	err = db.Get(user, sqlStr, user.UserName)

	if err != nil && err != sql.ErrNoRows {
		// 查询数据库出错
		return
	}
	if err == sql.ErrNoRows {
		// 用户不存在
		return ErrorUserNotExit
	}

	// 判断密码是否正确
	password := encryptPassword(oPassword)

	if password != user.Password {
		return ErrorPasswordWrong
	}

	return
}

/**
 * @Author ian-kevin
 * @Description //TODO 根据ID查询作者信息
 * @Date 22:05 2022/2/12
 **/
func GetUserByID(id int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	err = db.Get(user, sqlStr, id)
	return
}
