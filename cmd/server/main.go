package main

import (
	"log"
	"url-shortener/internal/config"
	"url-shortener/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	//load config
	defaultConfig, err := config.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	handler := handlers.NewURLHandler()

	router := gin.Default()
	router.POST("/urls", handler.CreateShortUrl)

	log.Println("Starting server on :8080")
	router.Run(":" + defaultConfig.Server.Port)
}
