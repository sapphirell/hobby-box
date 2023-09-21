package bootstrap

import (
	"github.com/go-redis/redis"
	"github.com/medivhzhan/weapp/v3"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisClient *redis.Client
var Wechat *weapp.Client
