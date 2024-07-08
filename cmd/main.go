package main

import (
	"github.com/Khvan-Group/blog-service/internal/api"
	blogApi "github.com/Khvan-Group/blog-service/internal/blogs/api"
	blogService "github.com/Khvan-Group/blog-service/internal/blogs/service"
	blogStore "github.com/Khvan-Group/blog-service/internal/blogs/store"
	commentApi "github.com/Khvan-Group/blog-service/internal/comments/api"
	commentService "github.com/Khvan-Group/blog-service/internal/comments/service"
	commentStore "github.com/Khvan-Group/blog-service/internal/comments/store"
	"github.com/Khvan-Group/blog-service/internal/db"
	"github.com/Khvan-Group/common-library/logger"
	"github.com/Khvan-Group/common-library/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"net/http"
)

const SERVER_PORT = "SERVER_PORT"

func main() {
	start()
}

func start() {
	// init logger
	logger.InitLogger()
	logger.Logger.Info("Starting server")

	// load environments
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// init db
	db.InitDB()

	port := ":" + utils.GetEnv(SERVER_PORT)
	newBlogStore := blogStore.New(db.DB)
	newBlogService := blogService.New(newBlogStore)
	newBlogApi := blogApi.New(*newBlogService, db.DB)

	newCommentStore := commentStore.New(db.DB)
	newCommentService := commentService.New(newCommentStore)
	newCommentApi := commentApi.New(*newCommentService, db.DB)

	srv := api.New(*newBlogApi, *newCommentApi, db.DB)
	r := mux.NewRouter()
	srv.AddRoutes(r)

	logger.Logger.Fatal(http.ListenAndServe(port, r).Error())
}
