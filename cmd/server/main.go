package main

import (
	"log"
	"time"

	// Import net/http for status codes
	api "url-shortener/generated" // Import the generated package
	"url-shortener/internal/config"
	"url-shortener/internal/db"
	"url-shortener/internal/handlers"
	"url-shortener/internal/repositories"
	"url-shortener/internal/services"
	"url-shortener/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
	defaultConfig, err := config.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := db.NewPostgresConnection(&defaultConfig.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	pgRepo := repositories.NewURLRepositoryPostgresql(dbConn)
	redisClient := db.NewRedisClient(&defaultConfig.Redis)
	redisRepo := repositories.NewURLRepositoryRedis(redisClient, time.Hour)
	timeProvider := utils.NewTimeProvider()
	urlRepo := repositories.NewURLRepository(redisRepo, pgRepo, timeProvider)
	urlStatPgRepo := repositories.NewURLStatisticsRepositoryPostgresql(dbConn)
	idGenerator := utils.NewNanoIDGenerator(12)

	urlService := services.NewURLService(urlRepo, urlStatPgRepo, idGenerator, timeProvider)
	urlStatService := services.NewURLStatsService(urlStatPgRepo)

	serverInterface := handlers.NewURLHandler(urlService, urlStatService, timeProvider)

	router := gin.New()
	router.Use(gin.LoggerWithFormatter(utils.CustomLogFormatter))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all domains (change for production)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	api.RegisterHandlersWithOptions(router, serverInterface, api.GinServerOptions{
		BaseURL: "", // Or set a base path if needed
		ErrorHandler: func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"message": err.Error()}) // Custom error handling
		},
	})

	log.Print("Starting server on :" + defaultConfig.Server.Port)
	router.Run(":" + defaultConfig.Server.Port)
}
