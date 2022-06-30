package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"ui-project/auth"
	"ui-project/lib"
	"ui-project/logger"
	ws "ui-project/websocket"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type usersDelivery struct {
	log         logger.LogUsecase
	userUsecase UserUsecase
	authUsecase auth.AuthUsecase
	ws          ws.Server
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewUserDelivery(ctx context.Context, ws ws.Server, log logger.LogUsecase, userUsecase UserUsecase, authUsecase auth.AuthUsecase) UserDelivery {
	return &usersDelivery{
		ws:          ws,
		log:         log,
		userUsecase: userUsecase,
		authUsecase: authUsecase,
	}
}

func (u *usersDelivery) Socket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := ws.Client{
		Id:   uuid.Must(uuid.NewRandom()).String(),
		Conn: conn,
	}
	u.ws.AddClient(client)

	// ... Use conn to send and receive messages.
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			u.log.LogErr(ctx, "socket client failed", err)
			u.ws.RemoveClient(client)
			return
		}
	}
}

// GetUserList
// @Summary  Get User list
// @Description Get User list
// @Tags User
// @Accept json
// @Produce json
// @param Authorization header string true "Authorization"
// @Param fullname query string false "fullname"
// @Param page query int false "page"
// @Param count query int false "count" Enums(10, 30, 50, 100)
// @Param orderBy query string false "orderBy" Enums(acct, pwd, fullname, create_at, update_at)
// @Param sort query string false "sort" Enums(asc, desc)
// @success 200 {array} User
// @Failure 400 {object} lib.ErrorResponse "Bad Request"
// @Security ApiKeyAuth
// @Router /users [GET]
func (u *usersDelivery) GetUserList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pagination := lib.NewPagination(1, 50, "acct", "asc")
	data, err := u.userUsecase.GetUserList(ctx, "", pagination)

	if err != nil {
		u.log.LogErr(ctx, "GetUserList failed", err)
		_, b := lib.ErrorResponseHelper(u.log.GetRequestId(ctx), "")
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

func (u *usersDelivery) GetUserListForQuery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	name := r.FormValue("fullname")
	page, pageErr := strconv.Atoi(r.FormValue("pagination"))
	count, countErr := strconv.Atoi(r.FormValue("count"))

	if pageErr != nil || countErr != nil {
		u.log.LogErr(ctx, "strconv.Atoi failed", pageErr, countErr)
		_, b := lib.ErrorResponseHelper(u.log.GetRequestId(ctx), "")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	}
	orderBy := r.FormValue("orderBy")
	sort := r.FormValue("sort")
	pagination := lib.NewPagination(page, count, orderBy, sort)

	if err := pagination.Verify([]string{"acct", "pwd", "fullname", "create_at", "update_at"}); err != nil {
		u.log.LogErr(ctx, "pagination verify failed", err)
		_, b := lib.ErrorResponseHelper(u.log.GetRequestId(ctx), err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	}
	data, err := u.userUsecase.GetUserList(ctx, name, pagination)

	if err != nil {
		u.log.LogErr(ctx, "GetUserList failed", err)
		_, b := lib.ErrorResponseHelper(u.log.GetRequestId(ctx), "")
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

// GetUser
// @Summary  Get User
// @Description Get User
// @Tags User
// @Accept json
// @Produce json
// @param Authorization header string true "Authorization"
// @success 200 {object} User
// @Failure 400 {object} lib.ErrorResponse "Bad Request"
// @Security ApiKeyAuth
// @Router /users/{account} [GET]
func (u *usersDelivery) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ctx := r.Context()
	data, err := u.userUsecase.GetUserByAccount(ctx, vars["account"])
	if err != nil {
		u.log.LogErr(ctx, "GetUser failed", err)
		_, b := lib.ErrorResponseHelper(u.log.GetRequestId(ctx), "")
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

// DeleteUser
// @Summary  Delete User
// @Description Delete User
// @Tags User
// @Accept json
// @Produce json
// @param Authorization header string true "Authorization"
// @success 200 {string} string "OK"
// @Failure 400 {object} lib.ErrorResponse "Bad Request"
// @Security ApiKeyAuth
// @Router /users/{account} [DELETE]
func (u *usersDelivery) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ctx := r.Context()
	err := u.userUsecase.DeleteUserByAccount(ctx, vars["account"])
	if err != nil {
		u.log.LogErr(ctx, "GetUser failed", err)
		_, b := lib.ErrorResponseHelper(u.log.GetRequestId(ctx), "")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("OK")
}

// UpdateUser
// @Summary  Update User
// @Description Update User
// @Tags User
// @Accept json
// @Produce json
// @param Authorization header string true "Authorization"
// @Param pwd body string false "pwd"
// @Param fullname body string false "fullname"
// @success 200 {string} string "OK"
// @Failure 400 {object} lib.ErrorResponse "Bad Request"
// @Security ApiKeyAuth
// @Router /users/{account} [PUT]
func (u *usersDelivery) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	user := User{}
	json.NewDecoder(r.Body).Decode(&user)

	err := u.userUsecase.UpdateUser(ctx, vars["account"], user)
	if err != nil {
		u.log.LogErr(ctx, "GetUser failed", err)
		_, b := lib.ErrorResponseHelper(u.log.GetRequestId(ctx), "")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("OK")
}

// CreateUsers
// @Summary  Create User
// @Description Create User
// @Tags User
// @Accept json
// @Produce json
// @Param acct body string true "acct"
// @Param pwd body string true "pwd"
// @Param fullname body string true "fullname"
// @success 200 {string} string "ok"
// @Failure 400 {object} lib.ErrorResponse "Bad Request"
// @Security ApiKeyAuth
// @Router /users [POST]
func (u *usersDelivery) CreateUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := User{}
	json.NewDecoder(r.Body).Decode(&user)
	if err := u.userUsecase.SetUser(ctx, user); err != nil {
		u.log.LogErr(ctx, "CreateUsers failed", err)
		_, b := lib.ErrorResponseHelper(u.log.GetRequestId(ctx), "")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("OK")
}

// Login
// @Summary  Login
// @Description Get jwt token
// @Tags User
// @Accept json
// @Produce json
// @Param acct body string true "acct"
// @Param pwd body string true "pwd"
// @success 200 {object} UserToken
// @Failure 400 {object} lib.ErrorResponse "Bad Request"
// @Security ApiKeyAuth
// @Router /users [POST]
func (u *usersDelivery) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := User{}
	json.NewDecoder(r.Body).Decode(&user)

	if check := u.userUsecase.VerifyUser(ctx, user); !check {
		u.log.LogErr(ctx, "Login failed")
		data, b := lib.ErrorResponseHelper(u.log.GetRequestId(ctx), "login failed")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		u.ws.Publish(data)
		return
	}

	token, err := u.authUsecase.GenerateToken(user.Acct)
	if err != nil {
		u.log.LogErr(ctx, "GenerateToken failed")
		data, b := lib.ErrorResponseHelper(u.log.GetRequestId(ctx), "login failed")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		u.ws.Publish(data)
		return
	}

	b, err := json.Marshal(UserToken{
		Token: token,
	})
	if err != nil {
		u.log.LogErr(ctx, "Marshal failed", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
