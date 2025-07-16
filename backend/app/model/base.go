package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID         uint      `gorm:"primaryKey;AUTO_INCREMENT" json:"id"`
	UUID       string    `gorm:"index:idx_uuid;unique;not null" json:"uuid"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.UUID = uuid.New().String()
	b.CreateTime = time.Now()
	b.UpdateTime = b.CreateTime
	return
}

func (b *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	b.UpdateTime = time.Now()
	return
}
