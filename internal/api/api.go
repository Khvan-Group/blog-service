package api

import (
	_ "github.com/Khvan-Group/blog-service/docs"
	blogApi "github.com/Khvan-Group/blog-service/internal/blogs/api"
	categoryApi "github.com/Khvan-Group/blog-service/internal/categories/api"
	commentApi "github.com/Khvan-Group/blog-service/internal/comments/api"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type API struct {
	blogApi     blogApi.API
	commentApi  commentApi.API
	categoryApi categoryApi.API
}

func New() *API {
	return &API{
		blogApi:     *blogApi.New(),
		commentApi:  *commentApi.New(),
		categoryApi: *categoryApi.New(),
	}
}

func (a *API) AddRoutes(r *mux.Router) {
	blogRouter := r.PathPrefix("/api/v1/blogs").Subrouter()
	a.blogApi.AddRoutes(blogRouter)

	commentRouter := r.PathPrefix("/api/v1/comments").Subrouter()
	a.commentApi.AddRoutes(commentRouter)

	categoryRouter := r.PathPrefix("/api/v1/categories").Subrouter()
	a.categoryApi.AddRoutes(categoryRouter)

	r.PathPrefix("/swagger").HandlerFunc(httpSwagger.WrapHandler)
}
