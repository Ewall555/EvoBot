package db

import (
	"EvoBot/backend/global"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func Init() {
	var dsn string
	switch global.CONF.DBConfig.Engine {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			global.CONF.DBConfig.Username,
			global.CONF.DBConfig.Password,
			global.CONF.DBConfig.Host,
			global.CONF.DBConfig.Port,
			global.CONF.DBConfig.Database,
		)
	case "postgres":
		dsn = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%d",
			global.CONF.DBConfig.Username,
			global.CONF.DBConfig.Password,
			global.CONF.DBConfig.Database,
			global.CONF.DBConfig.SSLMode,
			global.CONF.DBConfig.Host,
			global.CONF.DBConfig.Port,
		)
	default:
		global.ZAPLOG.Error("unsupported database engine: " + global.CONF.DBConfig.Engine)
		panic("unsupported database engine")
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	var db *gorm.DB
	var err error
	switch global.CONF.DBConfig.Engine {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "eb_",
				SingularTable: true,
			},
		})
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "eb_",
				SingularTable: true,
			},
		})
	default:
		global.ZAPLOG.Error("unsupported database engine: " + global.CONF.DBConfig.Engine)
		panic("unsupported database engine")
	}

	if err != nil {
		global.ZAPLOG.Error("failed to connect database, err:" + err.Error())
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		global.ZAPLOG.Error("connect db server failed, err:" + err.Error())
		panic(err)
	}

	sqlDB.SetConnMaxIdleTime(10 * time.Second)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	global.DB = db
	global.ZAPLOG.Info("init db successfully")
}
