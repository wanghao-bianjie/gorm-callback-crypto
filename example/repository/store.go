package repository

import (
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wanghao-bianjie/gorm-callback-crypto/example/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	dbLogger "gorm.io/gorm/logger"
)

var (
	db *gorm.DB
)

func InitMysqlDB(dsn string) {
	//db logger
	// logger
	newLogger := dbLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		dbLogger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      dbLogger.Warn, // Log level
			Colorful:      false,         // Disable color
		},
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		logrus.Fatalf(dsn)
	}
}

func GetDb() *gorm.DB {
	return db
}

func CreateTable() {
	if err := db.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(
		&model.User{},
	); err != nil {
		logrus.Error("AutoMigrate have error:", err.Error())
	}
}
