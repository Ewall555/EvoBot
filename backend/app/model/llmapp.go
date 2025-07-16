package model

type LLMAppConfig struct {
	BaseModel
	LLMAppType string `gorm:"column:llmapp_type;type:varchar(255);"`
	ConfigName string `gorm:"column:config_name;type:varchar(255);"`
	BaseURL    string `gorm:"column:base_url;type:varchar(255);"`
	ApiKey     string `gorm:"column:api_key;type:varchar(255);"`
}
