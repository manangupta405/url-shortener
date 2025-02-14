package handlers

import (
	"net/http"
	api "url-shortener/generated"

	"github.com/aidarkhanov/nanoid"
	"github.com/gin-gonic/gin"
)

type URLHandler struct {
}

func NewURLHandler() *URLHandler {
	return &URLHandler{}
}

func (h *URLHandler) CreateShortUrl(ctx *gin.Context) {
	var req api.CreateShortUrlJSONBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
		return
	}

	shortUrl := nanoid.New()

	response := &api.ShortenedUrlDetails{
		OriginalUrl: &req.OriginalUrl,
		ShortUrl:    &shortUrl}

	ctx.JSON(http.StatusCreated, response)
}
