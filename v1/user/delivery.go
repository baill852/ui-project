package user

import (
	"context"
	"encoding/json"
	"net/http"
	"ui-project/auth"
	"ui-project/lib"
	"ui-project/logger"

	"github.com/gorilla/mux"
)

type usersDelivery struct {
	log         logger.LogUsecase
	userUsecase UserUsecase
	authUsecase auth.AuthUsecase
}

func NewUserDelivery(ctx context.Context, log logger.LogUsecase, userUsecase UserUsecase, authUsecase auth.AuthUsecase) UserDelivery {
	return &usersDelivery{
		log:         log,
		userUsecase: userUsecase,
		authUsecase: authUsecase,
	}
}

func (u *usersDelivery) GetUserList(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("fullname")

	ctx := r.Context()
	data, err := u.userUsecase.GetUserList(ctx, name)
	if err != nil {
		u.log.LogErr(ctx, "GetUserList failed", err)
		b := lib.ErrorResponseHelper(u.log.GetRequestId(ctx), "")
		w.WriteHeader(http.StatusBadRequest)
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
		w.WriteHeader(http.StatusBadRequest)
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

func (u *usersDelivery) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ctx := r.Context()
	err := u.userUsecase.DeleteUserByAccount(ctx, vars["account"])
	if err != nil {
		u.log.LogErr(ctx, "GetUser failed", err)
		b := lib.ErrorResponseHelper(u.log.GetRequestId(ctx), "")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("OK")
}

func (u *usersDelivery) CreateUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := User{}
	json.NewDecoder(r.Body).Decode(&user)
	if err := u.userUsecase.SetUser(ctx, user); err != nil {
		u.log.LogErr(ctx, "CreateUsers failed", err)
		b := lib.ErrorResponseHelper(u.log.GetRequestId(ctx), "")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("OK")
}

func (u *usersDelivery) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := User{}
	json.NewDecoder(r.Body).Decode(&user)

	if check := u.userUsecase.VerifyUser(ctx, user); !check {
		u.log.LogErr(ctx, "Login failed")
		b := lib.ErrorResponseHelper(u.log.GetRequestId(ctx), "")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	}

	token, err := u.authUsecase.GenerateToken(user.Acct)
	if err != nil {
		u.log.LogErr(ctx, "GenerateToken failed")
		b := lib.ErrorResponseHelper(u.log.GetRequestId(ctx), "")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
