package repo

import (
	"gorm.io/gorm"
)

type DBOption func(*gorm.DB) *gorm.DB

type ICommonRepo interface {
	WithByID(id uint) DBOption
	WithByUUID(uuid string) DBOption
	WithIDsIn(ids []uint) DBOption
	WithUUIDsIn(uuid []string) DBOption
}

type CommonRepo struct{}

func NewCommonRepo() ICommonRepo {
	return &CommonRepo{}
}

func (c *CommonRepo) WithByID(id uint) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("id = ?", id)
	}
}

func (c *CommonRepo) WithByUUID(uuid string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("uuid = ?", uuid)
	}
}

func (c *CommonRepo) WithIDsIn(ids []uint) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("id in (?)", ids)
	}
}

func (c *CommonRepo) WithUUIDsIn(uuid []string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("uuid in (?)", uuid)
	}
}
