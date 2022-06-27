package user

import (
	"context"
	"ui-project/logger"

	"gorm.io/gorm"
)

type userRepository struct {
	client *gorm.DB
	log    logger.LogUsecase
}

func NewUserRepository(ctx context.Context, log logger.LogUsecase, client *gorm.DB) UserRepository {
	client.AutoMigrate(&User{})

	return &userRepository{
		client: client,
		log:    log,
	}
}

func (u *userRepository) GetUserList(ctx context.Context, name string) (data []User, err error) {
	result := u.client

	if len(name) > 0 {
		result = result.Where("fullname LIKE ?", "%"+name+"%")
	}

	result = result.Find(&data)
	return data, result.Error
}
