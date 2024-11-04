package database

import (
	"fmt"
	"github.com/nade-harlow/ecom-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"sync"
	"time"
)

var dbConnection *gorm.DB
var dbOnce sync.Once

func GetDbConnection() *gorm.DB {
	return dbConnection
}

func ConnectDb() *gorm.DB {
	dbOnce.Do(func() {
		dbConnection = setupDb()
	})

	return dbConnection
}

func setupDb() *gorm.DB {
	dbConfig := config.AppConfig.Database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Etc/UTC",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DatabaseName,
		dbConfig.Port)

	var logMod gormLogger.LogLevel = gormLogger.Warn
	if config.IsDev() {
		logMod = gormLogger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(logMod),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   dbConfig.Schema + ".",
			SingularTable: false,
		},
		SkipDefaultTransaction: true,
		NowFunc: func() time.Time {
			utcTz, err := time.LoadLocation("Africa/Lagos")
			if err != nil {
				panic("cannot get timezone for creating a now time. Error: " + err.Error())
			}
			return time.Now().In(utcTz)
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

	sqlDb, err := db.DB()
	if err != nil {
		panic("cannot get SQL Db, Reason: " + err.Error())
	}

	err = sqlDb.Ping()
	if err != nil {
		panic("cannot ping database, Reason: " + err.Error())
	}

	return db
}
