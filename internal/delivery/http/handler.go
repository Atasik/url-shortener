package delivery

import (
	v1 "link-shortener/internal/delivery/http/v1"
	"link-shortener/internal/service"
	"net/http"

	_ "link-shortener/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init() http.Handler {
	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	h.initAPI(r)

	return r
}

func (h *Handler) initAPI(router *mux.Router) {
	handlerV1 := v1.NewHandler(h.services)
	api := router.PathPrefix("/api").Subrouter()
	handlerV1.Init(api)
}
