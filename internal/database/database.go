package database

import "fmt"

func NewDatabase() (*gorm.DB, error) {
	fmt.Println("Setting up new database connection")
	return nil, nil
}