package main

import (
	"github.com/dkhvan-dev/alabs_project/blog-service/internal/blogs/service"
	"github.com/dkhvan-dev/alabs_project/blog-service/internal/blogs/store"
	"github.com/dkhvan-dev/alabs_project/blog-service/internal/db"
	"github.com/dkhvan-dev/alabs_project/common-libraries/logger"
	"github.com/dkhvan-dev/alabs_project/common-libraries/utils"
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
	r := mux.NewRouter()
	blogStore := store.New(db.DB)
	blogService := service.New(blogStore)

	logger.Logger.Fatal(http.ListenAndServe(port, r).Error())
}
