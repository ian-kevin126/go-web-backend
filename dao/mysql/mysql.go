package mysql

import (
	"fmt"
	"gin_demo/settings"

	"go.uber.org/zap"

	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Init(cfg *settings.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)

	// 也可以使用MustConnect连接，不成功就panic
	db, err = sqlx.Connect("mysql", dsn)

	if err != nil {
		zap.L().Error("Connect DB failed", zap.Error(err))
		return
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	return
}

// Close 对外暴露一个close数据库连接的方法
func Close() {
	_ = db.Close()
}
