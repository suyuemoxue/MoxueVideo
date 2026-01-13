package health

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Checker interface {
	Check(ctx context.Context) error
}

type MySQLChecker struct {
	db *gorm.DB
}

func NewMySQLChecker(db *gorm.DB) *MySQLChecker {
	return &MySQLChecker{db: db}
}

func (c *MySQLChecker) Check(ctx context.Context) error {
	if c == nil || c.db == nil {
		return errors.New("mysql unavailable")
	}
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
