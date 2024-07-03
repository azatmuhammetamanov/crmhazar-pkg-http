package crmhazar_pkg_http

import (
	"context"
	"github.com/gorilla/mux"
	"gitlab.com/salamtm.messenger/slog"
	"net/http"
)

type handler struct {
	ctx    context.Context
	logger *slog.Logger
}

func NotFoundHandler(ctx context.Context, logger *slog.Logger) Handler {
	return &handler{
		ctx:    ctx,
		logger: logger,
	}
}

func (h *handler) Register(router *mux.Router) {
	router.NotFoundHandler = http.HandlerFunc(h.NotFound)
}

func (h *handler) NotFound(w http.ResponseWriter, r *http.Request) {

	h.logger.Error("middleware not fund")

	w.WriteHeader(http.StatusNotFound)
	//w.Write("not Found")
}