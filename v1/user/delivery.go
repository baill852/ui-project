package user

import (
	"context"
	"encoding/json"
	"net/http"
)

type usersDelivery struct {
	userUsecase UserUsecase
}

func NewUserDelivery(ctx context.Context, userUsecase UserUsecase) UserDelivery {
	return &usersDelivery{
		userUsecase: userUsecase,
	}
}

func (u *usersDelivery) GetUserList(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
