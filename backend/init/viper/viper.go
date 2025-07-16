package viper

import (
	"EvoBot/backend/global"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init(configPath string) {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("viper.ReadInConfig failed, err: %s ", err))
	}

	if err := viper.Unmarshal(&global.CONF); err != nil {
		panic(fmt.Errorf("viper.Unmarshal failed, err: %s ", err))
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		global.ZAPLOG.Info("config file changed ...")
		if err := viper.Unmarshal(&global.CONF); err != nil {
			panic(fmt.Errorf("viper.Unmarshal failed, err: %s ", err))
		}
	})
}
