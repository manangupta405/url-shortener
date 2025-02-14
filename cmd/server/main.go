package main

import (
	"log"
	// Import net/http for status codes
	api "url-shortener/generated" // Import the generated package
	"url-shortener/internal/config"
	"url-shortener/internal/db"
	"url-shortener/internal/handlers"
	"url-shortener/internal/repositories"
	"url-shortener/internal/services"
	"url-shortener/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
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
	urlStatPgRepo := repositories.NewURLStatisticsRepositoryPostgresql(db)
	idGenerator := utils.NewNanoIDGenerator(6)
	timeProvider := utils.NewTimeProvider()
	urlService := services.NewURLService(pgRepo, urlStatPgRepo, idGenerator, timeProvider)
	urlStatService := services.NewURLStatsService(urlStatPgRepo)

	// Implement the generated ServerInterface:
	serverInterface := handlers.NewURLHandler(urlService, urlStatService)

	router := gin.Default()
	api.RegisterHandlersWithOptions(router, serverInterface, api.GinServerOptions{
		BaseURL: "", // Or set a base path if needed
		ErrorHandler: func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"message": err.Error()}) // Custom error handling
		},
	})

	log.Println("Starting server on :" + defaultConfig.Server.Port)
	router.Run(":" + defaultConfig.Server.Port)
}
