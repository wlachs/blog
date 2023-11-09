package repository

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectToMySQL establishes the DB connection
func connectToMySQL() (*gorm.DB, error) {
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
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
