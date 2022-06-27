package user

import (
	"context"

	"gorm.io/gorm"
)

type userRepository struct {
	client *gorm.DB
}

func NewUserRepository(ctx context.Context, client *gorm.DB) UserRepository {
	return &userRepository{
		client: client,
	}
}
