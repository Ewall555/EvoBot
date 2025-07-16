package client

import (
	"EvoBot/backend/global"
	"context"
	"strings"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/work"
	workConfig "github.com/silenceper/wechat/v2/work/config"
)

// var (
// 	wechatClient = initWechat()
// )

func initWechat() *wechat.Wechat {
	redisCache := cache.NewRedis(context.Background(), &cache.RedisOpts{})
	redisCache.SetConn(global.RDS)
	wechat := wechat.NewWechat()
	wechat.SetCache(redisCache)
	return wechat
}

func NewWork(conf WecomConfig) *work.Work {
	config := &workConfig.Config{
		CorpID:         conf.CorpID,
		CorpSecret:     conf.Secret,
		Token:          conf.Token,
		EncodingAESKey: conf.EncodingAESKey,
		AgentID:        conf.AgentID,
	}
	wechatClient := initWechat()
	return wechatClient.GetWork(config)
}

func splitContent(content string) []string {
	var parts []string
	currentPart := ""

	for _, word := range strings.Split(content, " ") {
		wordByteSize := len([]byte(word))

		if len(currentPart)+wordByteSize+1 > 2048 {
			parts = append(parts, strings.TrimSpace(currentPart))
			currentPart = ""
		}

		currentPart += word + " "
	}

	if currentPart != "" {
		parts = append(parts, strings.TrimSpace(currentPart))
	}

	if len(parts) > 5 {
		parts = parts[:5]
	}

	return parts
}
