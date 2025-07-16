package model

type KH struct {
	BaseModel
	KHID         string     `gorm:"uniqueIndex;column:khid;type:varchar(255);not null;" json:"khid"`
	StaffID      string     `gorm:"column:staffid;" json:"staffid"`
	VerifyStatus int        `gorm:"column:verify_status;type:int;not null;default:1;comment:'验证状态 1: 未处理 2: 角色确认 3：验证 4: 用户信息确认';" json:"verify_status"`
	ChatList     []ChatList `gorm:"foreignKey:KHID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ChatList struct {
	BaseModel
	KHID   uint   `gorm:"column:kh_id"`
	BotID  string `gorm:"column:botid;type:varchar(255);not null;" json:"botid"`
	ChatID string `gorm:"column:chatid;type:varchar(255);not null;" json:"chatid"`
}
