package repo

import (
	"EvoBot/backend/app/model"
	"EvoBot/backend/global"

	"gorm.io/gorm"
)

type KFRepo struct{}

type IKFRepo interface {
	Get(opts ...DBOption) (model.KF, error)
	WithKFID(kfid string) DBOption
	List(opts ...DBOption) ([]model.KF, error)
	Create(kf model.KF) error
	Update(kf model.KF) error
	Delete(opts ...DBOption) error
}

func NewIKFRepo() IKFRepo {
	return &KFRepo{}
}

// KF
func (r *KFRepo) Get(opts ...DBOption) (model.KF, error) {
	var kfinfo model.KF
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	db = db.Preload("Staffs.Policies.WorkTimes")
	if err := db.First(&kfinfo).Error; err != nil {
		return kfinfo, err
	}
	return kfinfo, nil
}

func (r *KFRepo) WithKFID(kfid string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("kfid = ?", kfid)
	}
}

func (r *KFRepo) List(opts ...DBOption) ([]model.KF, error) {
	var kf []model.KF
	db := global.DB.Model(&model.KF{})
	for _, opt := range opts {
		db = opt(db)
	}
	db = db.Preload("Staffs.Policies.WorkTimes")
	if err := db.Find(&kf).Error; err != nil {
		return kf, err
	}
	return kf, nil
}

func (r *KFRepo) Create(kf model.KF) error {
	return global.DB.Save(&kf).Error
}

func (r *KFRepo) Update(kf model.KF) error {
	if err := global.DB.Model(&kf).
		Where("id = ?", kf.ID).
		Updates(kf).Error; err != nil {
		return err
	}
	if err := global.DB.Model(&kf).Association("Staffs").Replace(kf.Staffs); err != nil {
		return err
	}
	return nil
}

func (r *KFRepo) Delete(opts ...DBOption) error {
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	return db.Delete(&model.KF{}).Error
}
