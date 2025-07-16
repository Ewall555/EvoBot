package migration

import (
	"EvoBot/backend/global"
	"EvoBot/backend/init/migration/migrations"

	"github.com/go-gormigrate/gormigrate/v2"
)

func Init() {
	m := gormigrate.New(global.DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migrations.AddTableUser,
		migrations.AddTableMaxkbConf,
		migrations.AddTableWecomConf,
		migrations.AddTableKH,
		migrations.AddTableChatList,
		migrations.AddTableKF,
		migrations.AddTableWorkTime,
		migrations.AddTablePolicy,
		migrations.AddTableStaff,
		migrations.AddTableKFStaff,
		migrations.AddTableStaffPolicy,
	})
	if err := m.Migrate(); err != nil {
		panic(err)
	}
	global.ZAPLOG.Info("Migration run successfully")
}
