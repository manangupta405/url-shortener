package main

import (
	"log"
	"url-shortener/internal/config"
	"url-shortener/internal/db"
	"url-shortener/internal/handlers"
	"url-shortener/internal/repositories"
	"url-shortener/internal/services"
	"url-shortener/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	//load config
	gin.SetMode(gin.ReleaseMode)
	defaultConfig, err := config.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := db.NewPostgresConnection(&defaultConfig.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	pgRepo := repositories.NewURLRepositoryPostgresql(db)
	idGenerator := utils.NewNanoIDGenerator(6)
	timeProvider := utils.NewTimeProvider()
	urlService := services.NewURLService(pgRepo, idGenerator, timeProvider)
	urlHandler := handlers.NewURLHandler(urlService)

	router := gin.Default()
	router.POST("/urls", urlHandler.CreateShortUrl)

	log.Println("Starting server on :" + defaultConfig.Server.Port)
	router.Run(":" + defaultConfig.Server.Port)
}
