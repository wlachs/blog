package db

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectToMySQL establishes the DB connection
func ConnectToMySQL() *gorm.DB {
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
		fmt.Println("failed to establish DB connection")
		os.Exit(1)
	}

	return db
}
