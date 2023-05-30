package models

import (
	"fmt"
	"time"

	"github.com/wlchs/blog/internal/database"
	"github.com/wlchs/blog/internal/transport/types"
)

type User struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	UserName     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Posts        []Post `gorm:"foreignKey:ID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func GetUsers() ([]User, error) {
	var u []User
	if result := database.Agent.Preload("Posts").Find(&u); result.Error != nil {
		return []User{}, result.Error
	}

	return u, nil
}

func GetUser(n string) (User, error) {
	u := User{
		UserName: n,
	}

	result := database.Agent.Preload("Posts").Take(&u)

	if result.Error != nil {
		return User{}, result.Error
	}

	if result.RowsAffected == 0 {
		return User{}, fmt.Errorf("user with name: %s not found", n)
	}

	return u, nil
}

func AddUser(u types.User) (User, error) {
	newUser := User{
		UserName:     u.UserName,
		PasswordHash: u.PasswordHash,
	}

	if result := database.Agent.Create(&newUser); result.Error != nil {
		return User{}, result.Error
	}

	return newUser, nil
}
