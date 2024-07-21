package main

import (
	_ "github.com/Khvan-Group/blog-service/docs"
	"github.com/Khvan-Group/blog-service/internal/api"
	"github.com/Khvan-Group/blog-service/internal/core/rabbitmq"
	"github.com/Khvan-Group/blog-service/internal/db"
	"github.com/Khvan-Group/common-library/logger"
	"github.com/Khvan-Group/common-library/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"net/http"
)

const SERVER_PORT = "SERVER_PORT"

// @title Blog Service API
// @version 1.0.3
// @description Blog Service.
// @host localhost:8082
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /
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

	// init RabbitMQ
	rabbitmq.InitRabbitMQ()

	// init db
	db.InitDB()

	port := ":" + utils.GetEnv(SERVER_PORT)
	r := mux.NewRouter()
	srv := api.New()
	srv.AddRoutes(r)

	logger.Logger.Fatal(http.ListenAndServe(port, r).Error())
}
