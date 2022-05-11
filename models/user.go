package models

type User struct {
	UserID   int64  `json:"user_id,string" db:"user_id"` // 指定json序列化/反序列化时使用小写user_id
	UserName string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Token    string
	//AccessToken  string
	//RefreshToken string
}

type RegisterForm struct {
	UserName        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}
