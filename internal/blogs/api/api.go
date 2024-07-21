package api

import (
	_ "github.com/Khvan-Group/blog-service/docs"
	"github.com/Khvan-Group/blog-service/internal/blogs/service"
	"github.com/Khvan-Group/common-library/constants"
	"github.com/Khvan-Group/common-library/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

type API struct {
	blogs service.Blogs
}

func New() *API {
	return &API{
		blogs: *service.New(),
	}
}

func (a *API) AddRoutes(r *mux.Router) {
	addSecurityRoutes(r, a)
	addUserRoutes(r, a)
}

func addSecurityRoutes(r *mux.Router, a *API) {
	r.Handle("", middlewares.AuthMiddleware(http.HandlerFunc(a.Create), constants.USER)).Methods(http.MethodPost)
	r.Handle("/{id}", middlewares.AuthMiddleware(http.HandlerFunc(a.Update), constants.USER)).Methods(http.MethodPut)
	r.Handle("/{id}", middlewares.AuthMiddleware(http.HandlerFunc(a.Delete))).Methods(http.MethodDelete)
	r.Handle("/{id}/send", middlewares.AuthMiddleware(http.HandlerFunc(a.Send))).Methods(http.MethodPost)
	r.Handle("/{id}", middlewares.AuthMiddleware(http.HandlerFunc(a.LikeOrFavorite))).Methods(http.MethodPost)
	r.Handle("/{username}/delete", middlewares.AuthMiddleware(http.HandlerFunc(a.DeleteAllByUsername), constants.ADMIN)).Methods(http.MethodDelete)
	r.Handle("/{id}/confirm", middlewares.AuthMiddleware(http.HandlerFunc(a.Confirm), constants.MODERATOR, constants.ADMIN)).Methods(http.MethodPost)
}

func addUserRoutes(r *mux.Router, a *API) {
	r.HandleFunc("", a.FindAll).Methods(http.MethodGet)
	r.HandleFunc("/{id}", a.FindById).Methods(http.MethodGet)
}
