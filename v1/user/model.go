package user

import (
	"context"
	"net/http"
	"time"
)

type User struct {
	Acct     string `gorm:"primaryKey"`
	Pwd      string
	FullName string
	CreateAt time.Time `gorm:"autoUpdateTime"`
	UpdateAt time.Time `gorm:"autoCreateTime"`
}

type UserDelivery interface {
	GetUserList(http.ResponseWriter, *http.Request)
}

type UserUsecase interface {
	GetUserList(context.Context) ([]User, error)
}

type UserRepository interface {
	GetUserList(context.Context) ([]User, error)
}
