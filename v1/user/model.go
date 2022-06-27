package user

import (
	"context"
	"net/http"
	"time"
)

type User struct {
	Acct     string `gorm:"primaryKey"`
	Pwd      string
	Fullname string
	CreateAt time.Time `gorm:"autoUpdateTime"`
	UpdateAt time.Time `gorm:"autoCreateTime"`
}

type UserDelivery interface {
	GetUserList(http.ResponseWriter, *http.Request)
}

type UserUsecase interface {
	GetUserList(ctx context.Context, name string) ([]User, error)
}

type UserRepository interface {
	GetUserList(ctx context.Context, name string) ([]User, error)
}
