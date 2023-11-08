package v1

import (
	"link-shortener/internal/service"

	"github.com/gorilla/mux"
)

const (
	appJSON = "application/json"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init(api *mux.Router) {
	r := api.PathPrefix("/v1").Subrouter()
	h.initLinkRoutes(r)
}
