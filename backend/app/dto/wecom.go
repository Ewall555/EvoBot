package dto

import wecomclient "EvoBot/backend/utils/wecom/client"

type SignatureOptions struct {
	wecomclient.SignatureOptions
}

type WecomConfig struct {
	Type           string `json:"type"`
	CorpID         string `json:"corp_id"`
	Token          string `json:"token"`
	EncodingAESKey string `json:"encoding_aes_key"`
	Secret         string `json:"secret"`
}
