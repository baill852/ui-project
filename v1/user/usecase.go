package user

import (
	"context"
)

type userUsecase struct {
	userRepository UserRepository
}

func NewUserUsecase(ctx context.Context, userRepository UserRepository) UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
	}
}
