package v1

import (
	"context"
	"ui-project/api"
	"ui-project/auth"
	"ui-project/logger"
	"ui-project/v1/user"
	ws "ui-project/websocket"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func Register(ctx context.Context, logger logger.LogUsecase, authUsecase auth.AuthUsecase, router *mux.Router, client *gorm.DB, ws ws.Server) []api.Route {
	// Repository
	userRepo := user.NewUserRepository(ctx, logger, client)

	// Usecase
	userUsecase := user.NewUserUsecase(ctx, logger, userRepo)

	// Delivery
	userDelivery := user.NewUserDelivery(ctx, ws, logger, userUsecase, authUsecase)

	return []api.Route{
		{
			Name:        "Socket",
			Method:      "GET",
			Pattern:     "/v1/socket",
			HandlerFunc: userDelivery.Socket,
			Secure:      false,
		},
		{
			Name:        "GetUserList",
			Method:      "GET",
			Pattern:     "/v1/users",
			HandlerFunc: userDelivery.GetUserList,
			Secure:      true,
		},
		// TODO: Regex not working
		{
			Name:    "GetUserList2",
			Method:  "GET",
			Pattern: "/v1/users",
			Queries: []string{
				"fullname", "{fullname}",
				"page", "{page:^[0-9]+$}",
				"count", "{count:^[0-9]+$}",
				"orderBy", "{orderBy}",
				"sort", "{sort}"},
			HandlerFunc: userDelivery.GetUserListForQuery,
			Secure:      true,
		},
		{
			Name:        "GetUser",
			Method:      "GET",
			Pattern:     "/v1/users/{account}",
			HandlerFunc: userDelivery.GetUser,
			Secure:      true,
		},
		{
			Name:        "DeleteUser",
			Method:      "DELETE",
			Pattern:     "/v1/users/{account}",
			HandlerFunc: userDelivery.DeleteUser,
			Secure:      true,
		},
		{
			Name:        "UpdateUser",
			Method:      "PUT",
			Pattern:     "/v1/users/{account}",
			HandlerFunc: userDelivery.UpdateUser,
			Secure:      true,
		},
		{
			Name:        "CreateUsers",
			Method:      "POST",
			Pattern:     "/v1/users",
			HandlerFunc: userDelivery.CreateUsers,
			Secure:      false,
		},
		{
			Name:        "Login",
			Method:      "POST",
			Pattern:     "/v1/users/login",
			HandlerFunc: userDelivery.Login,
			Secure:      false,
		},
	}
}
