package response

import "EvoBot/backend/app/dto"

type LLMAppConfig struct {
	UUID string `json:"uuid"`
	dto.LLMAppConfig
}
