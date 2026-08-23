package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase() *gorm.DB {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	idleConnection := GetEnvInt("DB_IDLE_CONNECTIONS", 10)
	maxConnection := GetEnvInt("DB_MAX_CONNECTIONS", 100)
	maxLifeTimeConnection := GetEnvInt("DB_MAX_LIFETIME_CONNECTIONS", 3600)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FJakarta",
		username,
		password,
		host,
		port,
		database,
	)

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database connection: %v", err)
	}

	sqlDB.SetMaxIdleConns(idleConnection)
	sqlDB.SetMaxOpenConns(maxConnection)
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifeTimeConnection) * time.Second)

	return db
}
