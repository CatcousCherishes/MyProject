package redis

import (
	"fmt"
	"web_app/settings"

	"github.com/go-redis/redis"
)

// go 连接redis的操作模板
// 声明一个全局的rdb变量

// var rdb *redis.Client
//
// // 初始化连接
//
//	func initClient() (err error) {
//		rdb = redis.NewClient(&redis.Options{
//			Addr:     "localhost:6379",
//			Password: "", // no password set
//			DB:       0,  // use default DB
//		})
//
//		_, err = rdb.Ping().Result()
//		if err != nil {
//			return err
//		}
//		return nil
//	}
//

//// 声明一个全局的rdb变量
//var rdb *redis.Client

var (
	client *redis.Client
	Nil    = redis.Nil
)

// Init   初始化连接
func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Password: cfg.Password, //  password set
		DB:       cfg.DB,       // use default DB
		PoolSize: cfg.PoolSize,
	})

	_, err = client.Ping().Result()
	return
}

func Close() {
	_ = client.Close()
}
