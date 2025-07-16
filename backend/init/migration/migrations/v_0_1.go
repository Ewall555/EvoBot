package migrations

import (
	"EvoBot/backend/app/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var AddTableMaxkbConf = &gormigrate.Migration{
	ID: "20250515-add-table-llmapp-config",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.LLMAppConfig{}); err != nil {
			return err
		}
		return nil
	},
}

var AddTableWecomConf = &gormigrate.Migration{
	ID: "20250515-add-table-wecom-config",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.WecomConfig{}); err != nil {
			return err
		}
		if err := tx.Create(&model.WecomConfig{
			Type: "addresslist",
		}).Error; err != nil {
			return err
		}
		if err := tx.Create(&model.WecomConfig{
			Type: "app",
		}).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddTableKH = &gormigrate.Migration{
	ID: "20250515-add-table-kh",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.KH{}); err != nil {
			return err
		}
		return nil
	},
}

var AddTableChatList = &gormigrate.Migration{
	ID: "20250515-add-table-chat-list",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.ChatList{}); err != nil {
			return err
		}
		return nil
	},
}

var AddTableKF = &gormigrate.Migration{
	ID: "20250515-add-table-kf",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.KF{}); err != nil {
			return err
		}
		return nil
	},
}

var AddTableKFStaff = &gormigrate.Migration{
	ID: "20250515-add-table-kf-staff",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.KF{}, &model.Staff{}); err != nil {
			return err
		}
		return nil
	},
}

var AddTableStaff = &gormigrate.Migration{
	ID: "20250515-add-table-staff",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Staff{}); err != nil {
			return err
		}
		return nil
	},
}

var AddTableStaffPolicy = &gormigrate.Migration{
	ID: "20250515-add-table-staff-policy",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Staff{}, &model.Policy{}); err != nil {
			return err
		}
		return nil
	},
}

var AddTablePolicy = &gormigrate.Migration{
	ID: "20250515-add-table-policy",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.Policy{}); err != nil {
			return err
		}
		return nil
	},
}

var AddTableWorkTime = &gormigrate.Migration{
	ID: "20250515-add-table-work-time",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.WorkTime{}); err != nil {
			return err
		}
		return nil
	},
}
