package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(DBUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(DBUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
