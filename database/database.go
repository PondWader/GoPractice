package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SessionCache struct {
	Username   string `gorm:"primaryKey"`
	Uuid       string
	Textures   string
	TimeCached time.Time
}

type DBConnOptions struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     uint16
}

func CreateDB(dbOpts *DBConnOptions) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dbOpts.User+":"+dbOpts.Password+"@tcp("+dbOpts.Host+":"+fmt.Sprint(dbOpts.Port)+")/"+dbOpts.Name+"?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&SessionCache{})

	go runDbJobs(db)

	return db
}
