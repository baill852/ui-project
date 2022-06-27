package database

import (
	"context"

	"gorm.io/gorm"
)

type DatabaseUsecase interface {
	GetClient() *gorm.DB
	Bootstrap(ctx context.Context)
	Shutdown(ctx context.Context)
}
