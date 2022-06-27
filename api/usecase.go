package api

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"ui-project/logger"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

const ContextRequestID string = "requestID"

type apiUsecase struct {
	http    *http.Server
	log     logger.LogUsecase
	handler *mux.Router
}

func NewApiUsecase(ctx context.Context, log logger.LogUsecase, handler *mux.Router) ApiUsecase {
	return &apiUsecase{
		log:     log,
		handler: handler,
	}
}

func (u *apiUsecase) RegisterRoute(ctx context.Context, routes []Route) {
	for _, route := range routes {
		var handler http.Handler
		if route.Secure {
			handler = u.authMiddleware(route.HandlerFunc)
		} else {
			handler = route.HandlerFunc
		}

		u.handler.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
}

func (u *apiUsecase) Bootstrap(ctx context.Context) {
	u.log.LogInfo(ctx, fmt.Sprintf("Listening and serving HTTP on %v:%v\n", viper.GetString("host"), viper.GetInt("port")))
	u.handler.NotFoundHandler = http.HandlerFunc(u.notFoundHandler)
	u.handler.MethodNotAllowedHandler = http.HandlerFunc(u.notAllowedHandler)
	u.handler.Use(u.loggerMiddleware)

	u.http = &http.Server{
		Addr: fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port")),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      u.applicationRecovery(u.headerMiddleware(u.handler)), // Pass our instance of gorilla/mux in.
	}
	go func() {
		if err := u.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			u.log.LogErr(ctx, "Mux Bootstrap failed", err)
			panic(err)
		}
	}()
}

func (u *apiUsecase) Shutdown(ctx context.Context) {
	if err := u.http.Shutdown(ctx); err != nil {
		u.log.LogErr(ctx, "Mux Shutdown failed", err)
	} else {
		u.log.LogInfo(ctx, "Mux Shutdown success")
	}
}

func (u *apiUsecase) applicationRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				u.log.LogErr(r.Context(), "Recovered from application error occurred", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (u *apiUsecase) headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (u *apiUsecase) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Add JWT authentication
		next.ServeHTTP(w, r)
	})
}

func (u *apiUsecase) loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		start := time.Now()
		ctx := u.log.SetRequestId(r.Context())
		r = r.WithContext(ctx)

		u.log.LogInfo(ctx, r.Method, r.RequestURI, time.Since(start).String())
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (u *apiUsecase) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func (u *apiUsecase) notAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}
