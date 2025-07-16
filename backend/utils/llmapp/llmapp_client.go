package llmapp

import (
	"EvoBot/backend/constant"
	"EvoBot/backend/utils/llmapp/client"
)

type LLMAppClient interface {
	ChatMessage(message string, chatid *string) (string, error)
	ChatOpen() (*string, error)
}

func NewLLMAppClient(llmappType string, vars map[string]interface{}) (LLMAppClient, error) {
	switch llmappType {
	case constant.Maxkb:
		return client.NewMaxkbClient(vars)
	default:
		return nil, constant.ErrNotSupportType
	}
}
