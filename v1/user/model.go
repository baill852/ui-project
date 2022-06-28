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
	GetUser(http.ResponseWriter, *http.Request)
	CreateUsers(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
}

type UserUsecase interface {
	GetUserByAccount(ctx context.Context, account string) (User, error)
	GetUserList(ctx context.Context, name string) ([]User, error)
	SetUser(ctx context.Context, user User) error
	VerifyUser(ctx context.Context, user User) bool
}

type UserRepository interface {
	GetUserByAccount(ctx context.Context, account string) (User, error)
	GetUserList(ctx context.Context, name string) ([]User, error)
	SetUser(ctx context.Context, user User) error
	VerifyUser(ctx context.Context, user User) bool
}
