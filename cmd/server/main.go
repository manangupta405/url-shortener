package main

import (
	"log"
	"url-shortener/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	handler := handlers.NewURLHandler()

	router := gin.Default()
	router.POST("/urls", handler.CreateShortUrl)

	log.Println("Starting server on :8080")
	router.Run(":8080")
}
