package connector

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"models"
)

type MysqlConnector struct {
}

func (m *MysqlConnector) getConnection() (db *gorm.DB, err error) {
	username := "lambda"
	password := "lambda"
	dbName := "test"
	host := "mysql"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, dbName)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func (m *MysqlConnector) GetDB() (db *gorm.DB, err error) {
	db, err = m.getConnection()
	if err != nil {
		return
	}

	err = db.AutoMigrate(&models.Unicorn{})
	if err != nil {
		return
	}

	return
}
