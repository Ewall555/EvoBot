package repo

import (
	"EvoBot/backend/app/model"
	"EvoBot/backend/global"
)

type PolicyRepo struct{}

type IPolicyRepo interface {
	Get(opts ...DBOption) (model.Policy, error)
	List(opts ...DBOption) ([]model.Policy, error)
	Create(policy model.Policy) error
	Update(policy model.Policy) error
	Delete(opts ...DBOption) error
	Page(page, size int, opts ...DBOption) (int64, []model.Policy, error)
	// CheckPolicyNameExists(name string) (bool, error)

	DeleteWorkTime(policyID uint) error
	CreateWorkTime(workTime model.WorkTime) error
}

func NewIPolicyRepo() IPolicyRepo {
	return &PolicyRepo{}
}

func (r *PolicyRepo) Get(opts ...DBOption) (model.Policy, error) {
	var policyinfo model.Policy

	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	db = db.Preload("WorkTimes")
	if err := db.First(&policyinfo).Error; err != nil {
		return policyinfo, err
	}

	return policyinfo, nil
}

func (r *PolicyRepo) Create(policy model.Policy) error {
	return global.DB.Create(&policy).Error
}

func (r *PolicyRepo) Update(policy model.Policy) error {
	return global.DB.Model(&model.Policy{}).
		Where("id = ?", policy.ID).
		Updates(policy).Error
}

func (r *PolicyRepo) List(opts ...DBOption) ([]model.Policy, error) {
	var policy []model.Policy
	db := global.DB.Model(&model.Policy{})
	for _, opt := range opts {
		db = opt(db)
	}
	db = db.Preload("WorkTimes")
	if err := db.Find(&policy).Error; err != nil {
		return policy, err
	}
	return policy, nil
}

func (r *PolicyRepo) Delete(opts ...DBOption) error {
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	return db.Delete(&model.Policy{}).Error
}

func (r *PolicyRepo) Page(page, size int, opts ...DBOption) (int64, []model.Policy, error) {
	var policies []model.Policy
	db := global.DB.Model(&model.Policy{})
	for _, opt := range opts {
		db = opt(db)
	}
	db = db.Preload("WorkTimes")
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, nil, err
	}
	if err := db.Limit(size).Offset(size * (page - 1)).Find(&policies).Error; err != nil {
		return count, nil, err
	}
	return count, policies, nil
}

// func (r *PolicyRepo) CheckPolicyNameExists(name string) (bool, error) {
// 	var count int64
// 	err := global.DB.Model(&model.Policy{}).
// 		Where("policyname = ?", name).
// 		Count(&count).Error
// 	return count > 0, err
// }

// WorkTime
func (r *PolicyRepo) DeleteWorkTime(policyID uint) error {
	return global.DB.Where("policy_id = ?", policyID).Delete(&model.WorkTime{}).Error
}

func (r *PolicyRepo) CreateWorkTime(workTime model.WorkTime) error {
	return global.DB.Create(&workTime).Error
}
