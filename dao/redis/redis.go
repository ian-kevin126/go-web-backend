package redis

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/go-redis/redis"
)

// 声明一个全局的rdb变量
var (
	client *redis.Client
	Nil    = redis.Nil
)

// Init 初始化连接
func Init() (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = client.Ping().Result()

	return err
}

func Close() {
	_ = client.Close()
}
