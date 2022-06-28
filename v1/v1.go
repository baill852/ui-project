package v1

import (
	"context"
	"ui-project/api"
	"ui-project/auth"
	"ui-project/logger"
	"ui-project/v1/user"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func Register(ctx context.Context, logger logger.LogUsecase, authUsecase auth.AuthUsecase, router *mux.Router, client *gorm.DB) []api.Route {
	// Repository
	userRepo := user.NewUserRepository(ctx, logger, client)

	// Usecase
	userUsecase := user.NewUserUsecase(ctx, logger, userRepo)

	// Delivery
	userDelivery := user.NewUserDelivery(ctx, logger, userUsecase, authUsecase)

	return []api.Route{
		{
			Name:        "GetUserList",
			Method:      "GET",
			Pattern:     "/v1/users",
			HandlerFunc: userDelivery.GetUserList,
			Secure:      true,
		},
		{
			Name:        "GetUserList",
			Method:      "GET",
			Pattern:     "/v1/users",
			Queries:     []string{"fullname", "{fullname}"},
			HandlerFunc: userDelivery.GetUserList,
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
