package postgres

import (
	"context"
	"fmt"

	"github.com/yehezkiel1086/go-redis-job-listings/internal/adapter/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func New(ctx context.Context, conf *config.DB) (*DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", conf.Host, conf.User, conf.Password, conf.Name, conf.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func (d *DB) GetDB() *gorm.DB {
	return d.db
}

func (d *DB) Migrate(dbs ...any) error {
	return d.db.AutoMigrate(dbs...)
}
