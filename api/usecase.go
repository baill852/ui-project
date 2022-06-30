package api

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"
	"ui-project/auth"
	"ui-project/logger"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

const ContextRequestID string = "requestID"

type apiUsecase struct {
	http    *http.Server
	log     logger.LogUsecase
	auth    auth.AuthUsecase
	handler *mux.Router
}

func NewApiUsecase(ctx context.Context, log logger.LogUsecase, auth auth.AuthUsecase, handler *mux.Router) ApiUsecase {
	return &apiUsecase{
		log:     log,
		auth:    auth,
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
			Handler(handler).
			Queries(route.Queries...).
			Name(route.Name)
	}
}

func (u *apiUsecase) Bootstrap(ctx context.Context) {
	u.log.LogInfo(ctx, fmt.Sprintf("Listening and serving HTTP on %v:%v\n", viper.GetString("host"), viper.GetInt("port")))
	u.handler.NotFoundHandler = http.HandlerFunc(u.notFoundHandler)
	u.handler.MethodNotAllowedHandler = http.HandlerFunc(u.notAllowedHandler)
	u.handler.Use(u.loggerMiddleware)
	u.handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := path.Dir("./index.html")
		w.Header().Set("Content-type", "text/html")
		http.ServeFile(w, r, p)
	})

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
		ctx := r.Context()
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := u.auth.ValidateToken(tokenString)
		if err != nil {
			u.log.LogErr(ctx, err)
			http.Error(w, "Error verifying JWT token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		account := claims.(jwt.MapClaims)["account"].(string)

		r.Header.Set("account", account)

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
