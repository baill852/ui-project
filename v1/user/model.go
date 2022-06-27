package user

import (
	"net/http"
)

type UserDelivery interface {
	GetUserList(http.ResponseWriter, *http.Request)
}

type UserUsecase interface {
}

type UserRepository interface {
}
