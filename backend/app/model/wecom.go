package model

type WecomConfig struct {
	BaseModel
	Type           string `gorm:"column:type;type:varchar(255);" json:"type"`
	CorpID         string `gorm:"column:corpid;type:varchar(255);" json:"corpid"`
	Token          string `gorm:"column:token;type:varchar(255);" json:"token"`
	EncodingAESKey string `gorm:"column:encoding_aeskey;type:varchar(255);" json:"encoding_aeskey"`
	Secret         string `json:"secret"`
	AgentID        string `json:"agent_id"`
}
