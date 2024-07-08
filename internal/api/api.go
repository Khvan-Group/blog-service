package api

import (
	blogApi "github.com/Khvan-Group/blog-service/internal/blogs/api"
	commentApi "github.com/Khvan-Group/blog-service/internal/comments/api"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	httpSwagger "github.com/swaggo/http-swagger"
)

type API struct {
	blogApi    blogApi.API
	commentApi commentApi.API
	DB         *sqlx.DB
}

func New(b blogApi.API, c commentApi.API, db *sqlx.DB) *API {
	return &API{
		blogApi:    b,
		commentApi: c,
		DB:         db,
	}
}

func (a *API) AddRoutes(r *mux.Router) {
	blogRouter := r.PathPrefix("/api/v1/blogs").Subrouter()
	a.blogApi.AddRoutes(blogRouter)
	
	commentRouter := r.PathPrefix("/api/v1/comments").Subrouter()
	a.commentApi.AddRoutes(commentRouter)
	
	r.PathPrefix("/swagger").HandlerFunc(httpSwagger.WrapHandler)
}
