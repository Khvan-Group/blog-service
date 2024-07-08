package api

import (
	"github.com/Khvan-Group/blog-service/internal/comments/service"
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
	comments service.Comments
	DB       *sqlx.DB
}

func New(c service.Comments, db *sqlx.DB) *API {
	return &API{
		comments: c,
		DB:       db,
	}
}

func (a *API) AddRoutes(r *mux.Router) {
	r.HandleFunc("/", a.FindAll).Methods(http.MethodGet)

	r.Handle("/", middlewares.AuthMiddleware(http.HandlerFunc(a.Create))).Methods(http.MethodPost)
	r.Handle("/{id}", middlewares.AuthMiddleware(http.HandlerFunc(a.Delete), MODERATOR, ADMIN)).Methods(http.MethodDelete)
}
