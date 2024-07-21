package api

import (
	"github.com/Khvan-Group/blog-service/internal/comments/service"
	"github.com/Khvan-Group/common-library/constants"
	"github.com/Khvan-Group/common-library/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

type API struct {
	comments service.Comments
}

func New() *API {
	return &API{
		comments: *service.New(),
	}
}

func (a *API) AddRoutes(r *mux.Router) {
	r.HandleFunc("", a.FindAll).Methods(http.MethodGet)

	r.Handle("", middlewares.AuthMiddleware(http.HandlerFunc(a.Create))).Methods(http.MethodPost)
	r.Handle("/{id}", middlewares.AuthMiddleware(http.HandlerFunc(a.Delete), constants.MODERATOR, constants.ADMIN)).Methods(http.MethodDelete)
}
