package crmhazar_pkg_http

import (
	"fmt"
	"net/http"
	"strconv"

	slog "github.com/azatmuhammetamanov/crmhazar-pkg-log"
)

type Middleware struct {
	logger  *slog.Logger
	jwtKey  string
	limiter *RateLimiter
}

type appBaseHandler func(w http.ResponseWriter, r *http.Request) Response
type appAuthHandler func(w http.ResponseWriter, r *http.Request, claims AuthClaims) Response

type AuthClaims struct {
	Id       int
	JobId    int
	BorderId int
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

		authClaims := AuthClaims{}

		id, err := strconv.Atoi(fmt.Sprint(claims["id"]))
		if err != nil {
			fmt.Println("2 err: ", err)

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		jobId, err := strconv.Atoi(fmt.Sprint(claims["job_id"]))
		if err != nil {
			fmt.Println("2 err: ", err)

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		borderId, err := strconv.Atoi(fmt.Sprint(claims["border_id"]))
		if err != nil {
			fmt.Println("2 err: ", err)

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authClaims.Id = id
		authClaims.JobId = jobId
		authClaims.BorderId = borderId

		result := h(w, r, authClaims)

		w.WriteHeader(result.GetStatusCode())
		w.Write(result.Marshal())
	}

}

func (middleware *Middleware) PAuth(h appAuthHandler) http.HandlerFunc {

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

		authClaims := AuthClaims{}
		fmt.Println(claims["id"])
		id, err := strconv.Atoi(fmt.Sprint(claims["id"]))
		if err != nil {
			fmt.Println("2 err: ", err)

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authClaims.Id = id

		result := h(w, r, authClaims)

		w.WriteHeader(result.GetStatusCode())
		w.Write(result.Marshal())
	}

}
