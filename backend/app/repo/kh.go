package repo

import (
	"EvoBot/backend/app/model"
	"EvoBot/backend/global"

	"gorm.io/gorm"
)

type KHRepo struct{}

type IKHRepo interface {
	Get(opts ...DBOption) (model.KH, error)
	WithKHID(khid string) DBOption
	Create(kh model.KH) error
	FirstOrCreatebyKHID(kh model.KH) error
	UpdatebyID(kh model.KH) error
	UpdatebyKHID(kh model.KH) error

	DeleteChatListWithBotID(botid string) error
	ClearChatIDsByKHIDAndBotID(khid string, botid string) error
	CreateChatList(chatlist model.ChatList) error
}

func NewIKHRepo() IKHRepo {
	return &KHRepo{}
}

func (r *KHRepo) Get(opts ...DBOption) (model.KH, error) {
	var kh model.KH
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	db = db.Preload("ChatList")
	if err := db.First(&kh).Error; err != nil {
		return kh, err
	}
	return kh, nil
}

func (r *KHRepo) WithKHID(khid string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("khid = ?", khid)
	}
}

func (r *KHRepo) Create(kh model.KH) error {
	return global.DB.Create(&kh).Error
}

func (r *KHRepo) FirstOrCreatebyKHID(kh model.KH) error {
	return global.DB.Where("khid = ?", kh.KHID).FirstOrCreate(&kh).Error
}

func (r *KHRepo) UpdatebyID(kh model.KH) error {
	return global.DB.Model(&model.KH{}).
		Where("id = ?", kh.ID).
		Updates(kh).Error
}

func (r *KHRepo) UpdatebyKHID(kh model.KH) error {
	return global.DB.Model(&model.KH{}).
		Where("khid = ?", kh.KHID).
		Updates(kh).Error
}

func (r *KHRepo) DeleteChatListWithBotID(botid string) error {
	return global.DB.Where("botid = ?", botid).Delete(&model.ChatList{}).Error
}

func (r *KHRepo) ClearChatIDsByKHIDAndBotID(khid string, botid string) error {
	var kh model.KH
	if err := global.DB.Where("khid = ?", khid).First(&kh).Error; err != nil {
		return err
	}
	result := global.DB.Where("kh_id = ? AND botid = ?", kh.ID, botid).Delete(&model.ChatList{})
	return result.Error
}

func (r *KHRepo) CreateChatList(chatlist model.ChatList) error {
	return global.DB.Create(&chatlist).Error
}
