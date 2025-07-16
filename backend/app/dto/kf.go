package dto

type KF struct {
	KFName           string `json:"kfname"`
	KFID             string `json:"kfid"`
	KFPlatform       string `json:"kfplatform"`
	BotID            string `json:"botid"`
	Status           int    `json:"status"`
	ReceivePriority  int    `json:"receive_priority"`
	ReceiveRule      int    `json:"receive_rule"`
	ChatTimeout      int    `json:"chat_timeout"`
	BotTimeout       int    `json:"bot_timeout"`
	BotTimeoutMsg    string `json:"bot_timeout_msg"`
	BotPrompt        string `json:"bot_prompt"`
	BotWelcomeMsg    string `json:"bot_welcome_msg"`
	StaffWelcomeMsg  string `json:"staff_welcome_msg"`
	UnmannedMsg      string `json:"unmanned_msg"`
	ChatendMsg       string `json:"chatend_msg"`
	TransferKeywords string `json:"transfer_keywords"`
}
