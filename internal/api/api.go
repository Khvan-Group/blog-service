package api

import (
	"database/sql"
	"github.com/dkhvan-dev/alabs_project/blog-service/internal/blogs/service"
	"github.com/dkhvan-dev/alabs_project/common-libraries/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

type API struct {
	blogs service.Blogs
	DB    *sql.DB
}

func New(b service.Blogs, db *sql.DB) *API {
	return &API{
		blogs: b,
		DB:    db,
	}
}

func (a *API) AddRoutes(r *mux.Router) {
	appRouter := r.PathPrefix("/api/v1/blogs").Subrouter()

	addUserRoutes(appRouter, a)
}

func addUserRoutes(r *mux.Router, a *API) {
	r.Handle("/", middlewares.AuthMiddleware("", http.HandlerFunc(a.Create))).Methods(http.MethodPost)
}
