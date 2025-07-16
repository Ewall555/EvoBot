package repo

import (
	"EvoBot/backend/app/model"
	"EvoBot/backend/global"

	"gorm.io/gorm"
)

type IStaffRepo interface {
	Get(opts ...DBOption) (model.Staff, error)
	WithStaffID(staffid string) DBOption
	WithStaffIdsIn(staffids []string) DBOption
	List(opts ...DBOption) ([]model.Staff, error)
	Create(staff model.Staff) error
	Update(staff model.Staff) error
	Delete(opts ...DBOption) error
}

type StaffRepo struct{}

func NewIStaffRepo() IStaffRepo {
	return &StaffRepo{}
}

// Staff
func (r *StaffRepo) Get(opts ...DBOption) (model.Staff, error) {
	var staffinfo model.Staff

	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	db = db.Preload("Policies.WorkTimes")
	if err := db.First(&staffinfo).Error; err != nil {
		return staffinfo, err
	}

	return staffinfo, nil
}

func (r *StaffRepo) WithStaffID(staffid string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("staffid = ?", staffid)
	}
}

func (c *StaffRepo) WithStaffIdsIn(staffids []string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("staffid in (?)", staffids)
	}
}

func (r *StaffRepo) List(opts ...DBOption) ([]model.Staff, error) {
	var staff []model.Staff
	db := global.DB.Model(&model.Staff{})
	for _, opt := range opts {
		db = opt(db)
	}
	db = db.Preload("Policies.WorkTimes")
	if err := db.Find(&staff).Error; err != nil {
		return staff, err
	}
	return staff, nil
}

func (r *StaffRepo) Create(staff model.Staff) error {
	return global.DB.Save(&staff).Error
}

func (r *StaffRepo) Update(staff model.Staff) error {
	if err := global.DB.Model(&model.Staff{}).
		Where("id = ?", staff.ID).
		Updates(staff).Error; err != nil {
		return err
	}
	if err := global.DB.Model(&staff).Association("Policies").Replace(staff.Policies); err != nil {
		return err
	}
	return nil
}

func (r *StaffRepo) Delete(opts ...DBOption) error {
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	return db.Delete(&model.Staff{}).Error
}
