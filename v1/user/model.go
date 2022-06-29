package user

import (
	"context"
	"net/http"
	"time"
	"ui-project/lib"
)

type User struct {
	Acct     string    `gorm:"primaryKey" example:"test"`
	Pwd      string    `example:"test"`
	Fullname string    `example:"test"`
	CreateAt time.Time `gorm:"autoUpdateTime" example:"2022-06-28T22:22:43.292795+08:00"`
	UpdateAt time.Time `gorm:"autoCreateTime" example:"2022-06-28T22:22:43.292795+08:00"`
}

type UserToken struct {
	Token string
}

type UserDelivery interface {
	GetUserList(http.ResponseWriter, *http.Request)
	GetUserListForQuery(w http.ResponseWriter, r *http.Request)
	GetUser(http.ResponseWriter, *http.Request)
	DeleteUser(http.ResponseWriter, *http.Request)
	UpdateUser(http.ResponseWriter, *http.Request)
	CreateUsers(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
}

type UserUsecase interface {
	GetUserByAccount(ctx context.Context, account string) (User, error)
	GetUserList(ctx context.Context, name string, pagination lib.Pagination) ([]User, error)
	SetUser(ctx context.Context, user User) error
	VerifyUser(ctx context.Context, user User) bool
	DeleteUserByAccount(ctx context.Context, account string) error
	UpdateUser(ctx context.Context, account string, user User) error
}

type UserRepository interface {
	GetUserByAccount(ctx context.Context, account string) (User, error)
	GetUserList(ctx context.Context, name string, pagination lib.Pagination) ([]User, error)
	SetUser(ctx context.Context, user User) error
	VerifyUser(ctx context.Context, user User) bool
	DeleteUserByAccount(ctx context.Context, account string) error
	UpdateUser(ctx context.Context, account string, user User) error
}
