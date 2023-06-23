package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Agent *gorm.DB
)

// InitDB establishes the DB connection
func InitDB() error {
	MysqlUser := os.Getenv("MYSQL_USER")
	MysqlPassword := os.Getenv("MYSQL_PASSWORD")
	MysqlDatabase := os.Getenv("MYSQL_DATABASE")
	MysqlHost := os.Getenv("MYSQL_HOST")
	MysqlPort := os.Getenv("MYSQL_PORT")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MysqlUser,
		MysqlPassword,
		MysqlHost,
		MysqlPort,
		MysqlDatabase,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	Agent = db

	return nil
}
