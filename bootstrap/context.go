package bootstrap

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisClient *redis.Client
