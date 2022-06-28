package user

import (
	"context"
	"ui-project/logger"

	"github.com/golang-jwt/jwt"
)

type authClaims struct {
	jwt.StandardClaims
	Account string `json:"acct"`
}

type userUsecase struct {
	log            logger.LogUsecase
	userRepository UserRepository
}

func NewUserUsecase(ctx context.Context, log logger.LogUsecase, userRepository UserRepository) UserUsecase {
	return &userUsecase{
		log:            log,
		userRepository: userRepository,
	}
}

func (u *userUsecase) VerifyUser(ctx context.Context, user User) bool {
	return u.userRepository.VerifyUser(ctx, user)
}

func (u *userUsecase) GetUserList(ctx context.Context, name string) ([]User, error) {
	return u.userRepository.GetUserList(ctx, name)
}

func (u *userUsecase) GetUserByAccount(ctx context.Context, account string) (User, error) {
	return u.userRepository.GetUserByAccount(ctx, account)
}

func (u *userUsecase) SetUser(ctx context.Context, user User) error {
	return u.userRepository.SetUser(ctx, user)
}

func (u *userUsecase) DeleteUserByAccount(ctx context.Context, account string) error {
	return u.userRepository.DeleteUserByAccount(ctx, account)
}
