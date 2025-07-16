package model

type User struct {
	BaseModel
	Username string `gorm:"index:idx_username;unique;type:varchar(256);not null" json:"username"`
	Password string `gorm:"type:varchar(256);not null" json:"-"`
}
