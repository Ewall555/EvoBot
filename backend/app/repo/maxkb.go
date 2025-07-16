package repo

import (
	"EvoBot/backend/app/model"
	"EvoBot/backend/global"
)

type ILLMAppRepo interface {
	Create(conf model.LLMAppConfig) error
	List(opts ...DBOption) ([]model.LLMAppConfig, error)
	Get(opts ...DBOption) (model.LLMAppConfig, error)
	Delete(opts ...DBOption) error
	Update(conf model.LLMAppConfig) error
}

type LLMAppRepo struct{}

func NewILLMAppRepo() ILLMAppRepo {
	return &LLMAppRepo{}
}

func (r *LLMAppRepo) List(opts ...DBOption) ([]model.LLMAppConfig, error) {
	var config []model.LLMAppConfig
	db := global.DB.Model(&model.LLMAppConfig{})
	for _, opt := range opts {
		db = opt(db)
	}
	if err := db.Find(&config).Error; err != nil {
		return config, err
	}
	return config, nil
}

func (r *LLMAppRepo) Get(opts ...DBOption) (model.LLMAppConfig, error) {
	var config model.LLMAppConfig

	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}

	if err := db.First(&config).Error; err != nil {
		return config, err
	}

	return config, nil
}

func (r *LLMAppRepo) Update(conf model.LLMAppConfig) error {
	return global.DB.Model(&model.LLMAppConfig{}).
		Where("id = ?", conf.ID).
		Updates(conf).Error
}

func (r *LLMAppRepo) Create(conf model.LLMAppConfig) error {
	return global.DB.Create(&conf).Error
}

func (r *LLMAppRepo) Delete(opts ...DBOption) error {
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	return db.Delete(&model.LLMAppConfig{}).Error
}
