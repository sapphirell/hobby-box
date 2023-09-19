package bootstrap

import (
	"github.com/go-redis/redis"
	"log"
	"os"
)

func InitRedis() {
	addr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	log.Printf("连接REDIS:%s; 使用密码%s", addr, os.Getenv("REDIS_PASS"))
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	pong, err := RedisClient.Ping().Result()
	if err != nil {
		log.Panicln("Redis无法链接", err)
	}
	log.Println("redis链接", pong, err)
}
