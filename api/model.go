package api

import (
	"context"
	"net/http"
)

type ApiUsecase interface {
	RegisterRoute(ctx context.Context, routes []Route)
	Bootstrap(ctx context.Context)
	Shutdown(ctx context.Context)
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	Queries     []string
	Secure      bool
	HandlerFunc http.HandlerFunc
}
