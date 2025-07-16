package dto

type LLMAppConfig struct {
	LLMAppType string `json:"llmapp_type"`
	ConfigName string `json:"config_name"`
	BaseURL    string `json:"base_url"`
	ApiKey     string `json:"api_key"`
}
