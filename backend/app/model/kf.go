package model

type KF struct {
	BaseModel
	KFName           string  `gorm:"column:kfname;type:varchar(255);" json:"kfname"`
	KFID             string  `gorm:"column:kfid;type:varchar(255);" json:"kfid"`
	KFPlatform       string  `gorm:"column:kfplatform;type:varchar(255);" json:"kfplatform"`
	BotID            string  `gorm:"column:botid;type:varchar(255);" json:"botid"`
	Status           int     `gorm:"column:status;type:int;comment:'1机器人可转人工,2仅机器人,3仅人工';" json:"status"`
	ReceivePriority  int     `gorm:"column:receive_priority;type:int;comment:'是否优先上一个';" json:"receive_priority"`
	ReceiveRule      int     `gorm:"column:receive_rule;type:int;comment:'1轮流接待，2空闲接待';" json:"receive_rule"`
	ChatTimeout      int     `gorm:"column:chat_timeout;comment:'会话超时时间';" json:"chat_timeout"`
	BotTimeout       int     `gorm:"column:bot_timeout;comment:'机器人超时时间';" json:"bot_timeout"`
	BotTimeoutMsg    string  `gorm:"column:bot_timeout_msg;type:varchar(255);comment:'机器人超时提示语';" json:"bot_timeout_msg"`
	BotWelcomeMsg    string  `gorm:"column:bot_welcome_msg;type:varchar(255);comment:'机器人的欢迎语';" json:"botwelcome_msg"`
	StaffWelcomeMsg  string  `gorm:"column:staff_welcome_msg;type:varchar(255);comment:'转人工的欢迎语';" json:"staff_welcome_msg"`
	UnmannedMsg      string  `gorm:"column:unmanned_msg;type:varchar(255);comment:'无人值守语';" json:"unmanned_msg"`
	ChatendMsg       string  `gorm:"column:chatend_msg;type:varchar(255);comment:'会话结束语';" json:"chatend_msg"`
	TransferKeywords string  `gorm:"column:transfer_keywords;type:varchar(255);comment:'转人工关键字列表';" json:"transfer_keywords"`
	Staffs           []Staff `gorm:"many2many:kf_staff;foreignKey:ID;joinForeignKey:KFID;References:ID;joinReferences:StaffID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"kf_staffs"`
}
