package model

type Policy struct {
	BaseModel
	PolicyName string     `gorm:"column:policyname;type:varchar(255);unique;" json:"policyname"`
	Repeat     int        `gorm:"column:repeat;type:int;" json:"repeat"`       // 1.自定义[周一至周天] 2.每天 3.周一至周五 4.法定工作日(跳过法定节假日) 5.法定节假日(跳过法定工作日)
	Week       string     `gorm:"column:week;type:varchar(7);" json:"week"`    // Repeat为1时必填，"1111100" 从左往右周日、周一、周二、周三、周四、周五、周六
	MaxCount   int        `gorm:"column:max_count;type:int;" json:"max_count"` // 每人最大接待数量
	WorkTimes  []WorkTime `gorm:"foreignKey:PolicyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type WorkTime struct {
	BaseModel
	PolicyID  uint   `gorm:"column:policy_id"`
	StartTime string `gorm:"column:start_time;type:time(0);" json:"start_time"` // HH:MM:SS
	EndTime   string `gorm:"column:end_time;type:time(0);" json:"end_time"`     // HH:MM:SS
}
