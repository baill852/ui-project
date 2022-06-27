package v1

import (
	"context"
	"ui-project/api"
	"ui-project/v1/user"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func Register(ctx context.Context, router *mux.Router, client *gorm.DB) []api.Route {
	// Repository

	userRepo := user.NewUserRepository(ctx, client)

	// Usecase
	userUsecase := user.NewUserUsecase(ctx, userRepo)

	// Delivery
	userDelivery := user.NewUserDelivery(ctx, userUsecase)

	return []api.Route{
	}
}
