package repo

import (
	"EvoBot/backend/app/model"
	"EvoBot/backend/global"

	"gorm.io/gorm"
)

type AuthRepo struct{}

type IAuthRepo interface {
	Get(opts ...DBOption) (model.User, error)
	WithByUsername(Username string) DBOption
}

func NewISettingRepo() IAuthRepo {
	return &AuthRepo{}
}

func (u *AuthRepo) Get(opts ...DBOption) (model.User, error) {
	var user model.User
	db := global.DB.Model(&model.User{})
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.First(&user).Error
	return user, err
}

func (c *AuthRepo) WithByUsername(Username string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("Username = ?", Username)
	}
}
