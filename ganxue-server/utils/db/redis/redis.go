package redis

import (
	"ganxue-server/global"

	"github.com/gangantongxue/ggl"
	"github.com/go-redis/redis/v8"
)

func Init() {
	global.RDB = redis.NewClient(&redis.Options{
		Addr:     global.CONFIG.Redis.Addr,
		Password: global.CONFIG.Redis.Password,
		DB:       global.CONFIG.Redis.DB,
	})
	// 测试连接
	_, err := global.RDB.Ping(global.CTX).Result()
	if err != nil {
		ggl.Fatal("redis连接失败", ggl.Err(err))
	}
	ggl.Info("redis连接成功")
}

func Close() {
	err := global.RDB.Close()
	if err != nil {
		ggl.Fatal("redis关闭失败", ggl.Err(err))
	}
}
