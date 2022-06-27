package user

import (
	"context"
	"encoding/json"
	"net/http"
	"ui-project/logger"
)

type usersDelivery struct {
	log         logger.LogUsecase
	userUsecase UserUsecase
}

func NewUserDelivery(ctx context.Context, log logger.LogUsecase, userUsecase UserUsecase) UserDelivery {
	return &usersDelivery{
		log:         log,
		userUsecase: userUsecase,
	}
}

func (u *usersDelivery) GetUserList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data, err := u.userUsecase.GetUserList(ctx)
	if err != nil {
		u.log.LogErr(ctx, "GetUserList failed", err)
	}

	b, err := json.Marshal(data)
	if err != nil {
		u.log.LogErr(ctx, "Marshal failed", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
