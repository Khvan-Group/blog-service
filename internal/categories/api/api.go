package api

import (
	"github.com/Khvan-Group/blog-service/internal/categories/service"
	"github.com/Khvan-Group/common-library/constants"
	"github.com/Khvan-Group/common-library/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

type API struct {
	categories service.Categories
}

func New() *API {
	return &API{
		categories: *service.New(),
	}
}

func (a *API) AddRoutes(r *mux.Router) {
	addRoutes(r, a)
	addSecurityRoutes(r, a)
}

func addRoutes(r *mux.Router, a *API) {
	r.HandleFunc("", a.FindAll).Methods(http.MethodGet)
}

func addSecurityRoutes(r *mux.Router, a *API) {
	r.Handle("", middlewares.AuthMiddleware(http.HandlerFunc(a.Save), constants.MODERATOR, constants.ADMIN)).Methods(http.MethodPost)
	r.Handle("/{code}", middlewares.AuthMiddleware(http.HandlerFunc(a.Delete), constants.MODERATOR, constants.ADMIN)).Methods(http.MethodDelete)
}
