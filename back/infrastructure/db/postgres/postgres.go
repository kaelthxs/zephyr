package postgres

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "zephyr-backend/config"
)

func NewPostgres(cfg *config.Config) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return db, nil
}
