package user

import (
	"context"
	"ui-project/lib"
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

func (u *userRepository) GetUserByAccount(ctx context.Context, account string) (data User, err error) {
	result := u.client.First(&data, "acct = ?", account)

	return data, result.Error
}

func (u *userRepository) GetUserList(ctx context.Context, name string, pagination lib.Pagination) (data []User, err error) {
	result := u.client

	if len(name) > 0 {
		result = result.Where("fullname LIKE ?", "%"+name+"%")
	}

	result = result.Offset(pagination.GetOffset()).Limit(pagination.GetCount()).Order(pagination.GetOrderString()).Find(&data)
	return data, result.Error
}

func (u *userRepository) SetUser(ctx context.Context, user User) error {
	result := u.client.Create(&user)
	return result.Error
}

func (u *userRepository) VerifyUser(ctx context.Context, user User) bool {
	result := u.client.First(&user, "acct = ? AND pwd = ?", user.Acct, user.Pwd)

	if result.Error != nil {
		u.log.LogErr(ctx, "VerifyUser", result.Error)
		return false
	}

	if result.RowsAffected == 0 {
		return false
	}

	return true
}

func (u *userRepository) DeleteUserByAccount(ctx context.Context, account string) error {
	result := u.client.Where("acct = ?", account).Delete(&User{})
	return result.Error
}

func (u *userRepository) UpdateUser(ctx context.Context, account string, user User) error {
	result := u.client.Model(&User{}).
		Omit("acct", "create_at", "update_at").
		Where("acct = ?", account).
		Updates(user)

	return result.Error
}
