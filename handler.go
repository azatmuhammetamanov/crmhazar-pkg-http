package crmhazar_pkg_http

import (
	"github.com/gorilla/mux"
)

type Handler interface {
	Register(router *mux.Router)
}
