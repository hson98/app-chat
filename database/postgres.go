package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"hson98/app-chat/config"
)

func NewPostgresDB(config *config.Config) *gorm.DB {
	username := config.DBUser
	password := config.DBPass
	dbName := config.DBName
	dbHost := config.DBHost
	port := config.DBPort
	//dsn = "host=localhost auth=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"require
	dsn := fmt.Sprintf("host=%s auth=%s password=%s dbname=%s port=%s sslmode=disable  TimeZone=Asia/Ho_Chi_Minh", dbHost, username, password, dbName, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return db
}
