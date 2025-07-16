package global

import (
	"EvoBot/backend/configs"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	RDS    redis.UniversalClient
	ZAPLOG *zap.Logger
	CONF   configs.ServerConfig
)
