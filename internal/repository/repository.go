package repository

import (
	"gorm.io/gorm"
)

// Repository defines the database access layer
type Repository interface {
	Select(query interface{}, args ...interface{}) *gorm.DB
	Find(out interface{}, where ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Updates(value interface{}) *gorm.DB
	Delete(value interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Preload(column string, conditions ...interface{}) *gorm.DB
	Close() error
	AutoMigrate(value interface{}) error
}

// repository implements the Repository interface and stores the concrete Gorm DB implementation
type repository struct {
	db *gorm.DB
}

// CreateRepository established a DB connection and returns the repository.
// If, for some reason, the DB connection fails, the error is logged and the application terminates.
func CreateRepository(database *gorm.DB) Repository {
	return &repository{db: database}
}

// Select specify fields to be retrieved from the database
func (rep *repository) Select(query interface{}, args ...interface{}) *gorm.DB {
	return rep.db.Select(query, args...)
}

// Find retrieves rows satisfying the given conditions
func (rep *repository) Find(out interface{}, where ...interface{}) *gorm.DB {
	return rep.db.Find(out, where...)
}

// Create inserts a value in the database
func (rep *repository) Create(value interface{}) *gorm.DB {
	return rep.db.Create(value)
}

// Updates selected rows in the database
func (rep *repository) Updates(value interface{}) *gorm.DB {
	return rep.db.Updates(value)
}

// Delete removes rows from the database
func (rep *repository) Delete(value interface{}) *gorm.DB {
	return rep.db.Delete(value)
}

// Where filters the retrieved objects from the database
func (rep *repository) Where(query interface{}, args ...interface{}) *gorm.DB {
	return rep.db.Where(query, args...)
}

// Preload loads a foreign table to run queries based on joined tables
func (rep *repository) Preload(column string, conditions ...interface{}) *gorm.DB {
	return rep.db.Preload(column, conditions...)
}

// Close closes the database connection
func (rep *repository) Close() error {
	sqlDB, _ := rep.db.DB()
	return sqlDB.Close()
}

// AutoMigrate updates the DB schema to match the current state
func (rep *repository) AutoMigrate(value interface{}) error {
	return rep.db.AutoMigrate(value)
}
