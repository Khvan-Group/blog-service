package api

import (
	"github.com/Khvan-Group/blog-service/internal/blogs/service"
	"github.com/Khvan-Group/common-library/middlewares"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
)

const (
	ADMIN     = "ADMIN"
	MODERATOR = "MODERATOR"
	USER      = "USER"
)

type API struct {
	blogs service.Blogs
	DB    *sqlx.DB
}

func New(b service.Blogs, db *sqlx.DB) *API {
	return &API{
		blogs: b,
		DB:    db,
	}
}

func (a *API) AddRoutes(r *mux.Router) {
	addSecurityRoutes(r, a)
	addUserRoutes(r, a)
}

func addSecurityRoutes(r *mux.Router, a *API) {
	r.Handle("/", middlewares.AuthMiddleware(http.HandlerFunc(a.Create), USER)).Methods(http.MethodPost)
	r.Handle("/", middlewares.AuthMiddleware(http.HandlerFunc(a.Update), USER)).Methods(http.MethodPut)
	r.Handle("/{id}", middlewares.AuthMiddleware(http.HandlerFunc(a.Delete))).Methods(http.MethodDelete)
	r.Handle("/{id}/{action}", middlewares.AuthMiddleware(http.HandlerFunc(a.LikeOrFavorite))).Methods(http.MethodPost)
	r.Handle("/{id}/confirm", middlewares.AuthMiddleware(http.HandlerFunc(a.Confirm), MODERATOR, ADMIN)).Methods(http.MethodPost)
}

func addUserRoutes(r *mux.Router, a *API) {
	r.HandleFunc("/", a.FindAll).Methods(http.MethodGet)
	r.HandleFunc("/{id}", a.FindById).Methods(http.MethodGet)
}
