package client

import (
	"EvoBot/backend/global"

	"github.com/Ewall555/MaxKB-golang-sdk/api/application"
	mkreq "github.com/Ewall555/MaxKB-golang-sdk/api/request"   // 请求参数
	mkresp "github.com/Ewall555/MaxKB-golang-sdk/api/response" // 返回参数
	mkconfig "github.com/Ewall555/MaxKB-golang-sdk/config"     // 配置参数
	mk "github.com/Ewall555/MaxKB-golang-sdk/maxkb"            // 引入包
	"go.uber.org/zap"
)

type MaxKB struct {
	MaxKBClient     *mk.MaxKB
	ApplicationChat *application.ApplicationChat
}

func NewMaxkbClient(vars map[string]interface{}) (*MaxKB, error) {
	base_url := loadParamFromVars("base_url", vars)
	api_key := loadParamFromVars("api_key", vars)
	maxKBClient := mk.NewMaxKB(&mkconfig.Config{
		BaseURL: base_url,
		ApiKey:  api_key,
	})
	applicationChat := maxKBClient.GetApplicationChat()
	return &MaxKB{
		MaxKBClient:     maxKBClient,
		ApplicationChat: applicationChat,
	}, nil
}

type StreamCallback struct {
	StreamCallback func(*mkresp.Chat_messagePostStreamResponse)
}

func (c *MaxKB) ChatMessage(message string, chatid *string) (string, error) {
	req := mkreq.Chat_messagePostRequest{
		Message: message,
		ReChat:  false,
		Stream:  false,
	}
	resp, err := c.ApplicationChat.Chat_messageByChat_id(req, chatid, nil)
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}

func (c *MaxKB) ChatOpen() (*string, error) {
	profileresp, err := c.ApplicationChat.Profile()
	if err != nil {
		global.ZAPLOG.Error("get profile error", zap.Error(err))
		return nil, err
	}
	chatid, err := c.ApplicationChat.ChatOpenByApplication_id(profileresp.ID)
	if err != nil {
		global.ZAPLOG.Error("open chat error", zap.Error(err))
		return nil, err
	}
	return chatid, err
}
