package crmhazar_pkg_http

import (
	"fmt"
	"net/http"
	"strconv"

	slog "gitlab.com/GadamGurbanov/crmhazar-pkg-log"
)

type Middleware struct {
	logger  *slog.Logger
	jwtKey  string
	limiter *RateLimiter
}

type appBaseHandler func(w http.ResponseWriter, r *http.Request) Response
type appAuthHandler func(w http.ResponseWriter, r *http.Request, claims AuthClaims) Response

type AuthClaims struct {
	Id         int64
	DeviceId   int64
	AppVersion string
}

func NewMiddleware(logger *slog.Logger, jwtKey string, limiter *RateLimiter) *Middleware {
	return &Middleware{
		logger:  logger,
		jwtKey:  jwtKey,
		limiter: limiter,
	}
}

func (middleware *Middleware) Base(h appBaseHandler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		result := h(w, r)

		w.WriteHeader(result.GetStatusCode())
		w.Write(result.Marshal())
	}
}

func (middleware *Middleware) Auth(h appAuthHandler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		token := r.Header.Get("authorization")

		claims, err := TokenClaims(token, middleware.jwtKey)

		if err != nil {

			if err.Error() == "Token is expired" {
				w.WriteHeader(http.StatusNotAcceptable)
			}

			middleware.logger.Error("shttp error jwt: ", err)

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authClaims := AuthClaims{
			AppVersion: fmt.Sprint(claims["app_version"]),
		}

		id, err := strconv.ParseInt(fmt.Sprint(claims["user_id"]), 10, 64)
		if err != nil {
			fmt.Println("2 err: ", err)

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		deviceId, err := strconv.ParseInt(fmt.Sprint(claims["device_id"]), 10, 64)
		if err != nil {
			fmt.Println("3 err: ", err)

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		//allowed := middleware.limiter.Allow(fmt.Sprint(deviceId))
		//
		//if !allowed {
		//	fmt.Println("from deviceID: ", deviceId, " blocked")
		//	w.WriteHeader(http.StatusTooManyRequests)
		//	return
		//}

		authClaims.Id = id
		authClaims.DeviceId = deviceId

		result := h(w, r, authClaims)

		w.WriteHeader(result.GetStatusCode())
		w.Write(result.Marshal())
	}

}
