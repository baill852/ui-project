package user

import (
	"context"
	"encoding/json"
	"net/http"
	"ui-project/lib"
	"ui-project/logger"

	"github.com/gorilla/mux"
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
	name := r.FormValue("fullname")

	ctx := r.Context()
	data, err := u.userUsecase.GetUserList(ctx, name)
	if err != nil {
		u.log.LogErr(ctx, "GetUserList failed", err)
		b := lib.ErrorResponseHelper(u.log.GetRequestId(ctx), "")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(b)
		return
	}

	b, err := json.Marshal(data)
	if err != nil {
		u.log.LogErr(ctx, "Marshal failed", err)

	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (u *usersDelivery) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ctx := r.Context()
	data, err := u.userUsecase.GetUserByAccount(ctx, vars["account"])
	if err != nil {
		u.log.LogErr(ctx, "GetUser failed", err)
		b := lib.ErrorResponseHelper(u.log.GetRequestId(ctx), "")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(b)
		return
	}

	b, err := json.Marshal(data)
	if err != nil {
		u.log.LogErr(ctx, "Marshal failed", err)

	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
