package server

import (
	"EvoBot/backend/global"
	"EvoBot/backend/init/db"
	"EvoBot/backend/init/migration"
	"EvoBot/backend/init/redis"
	"EvoBot/backend/init/router"
	"EvoBot/backend/init/viper"
	"EvoBot/backend/init/zaplog"
	"fmt"
)

func Start(configPath string) {
	viper.Init(configPath)
	zaplog.Init()

	redis.Init()
	db.Init()
	migration.Init()

	r := router.Init()

	r.Run(fmt.Sprintf("%s:%d", global.CONF.SystemConfig.Host, global.CONF.SystemConfig.Port))
}
