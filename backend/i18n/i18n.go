package i18n

import (
	"EvoBot/backend/global"
	"EvoBot/backend/i18n/lang"
	"fmt"
)

func T(key string) string {
	switch global.CONF.LogConfig.Lang {
	case "zh":
		if val, ok := lang.ZhMessages[key]; ok {
			return val
		}
	case "en":
		if val, ok := lang.EnMessages[key]; ok {
			return val
		}
	}
	return key
}

func Tf(key string, args ...any) string {
	return fmt.Sprintf(T(key), args...)
}
