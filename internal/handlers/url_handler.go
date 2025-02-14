package handlers

import (
	"log"
	"net/http"
	api "url-shortener/generated"
	"url-shortener/internal/services"

	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	service        services.URLService
	urlStatService services.URLStatsService
	api.ServerInterface
}

func NewURLHandler(service services.URLService, urlStatService services.URLStatsService) *URLHandler {
	return &URLHandler{service: service, urlStatService: urlStatService}
}

func (h *URLHandler) CreateShortUrl(ctx *gin.Context) {
	var req api.CreateShortUrlJSONBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
		return
	}

	shortPath, err := h.service.CreateShortURL(ctx, req.OriginalUrl, req.Expiry)
	if err != nil {
		log.Printf(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create short URL"})
		return
	}

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	shortUrl := scheme + "://" + ctx.Request.Host + "/" + shortPath

	response := &api.ShortenedUrlDetails{
		OriginalUrl: &req.OriginalUrl,
		ShortUrl:    &shortUrl}

	ctx.JSON(http.StatusCreated, response)
}

func (h *URLHandler) RedirectToOriginalUrl(ctx *gin.Context, shortPath string) {
	longURL, err := h.service.GetLongURL(ctx, shortPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to redirect"})
		return
	}

	if longURL == "" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, longURL)
}

// DeleteShortURL implements URLService
func (h *URLHandler) DeleteShortUrl(ctx *gin.Context, shortPath string) {
	err := h.service.DeleteURL(ctx, shortPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": shortPath + " deleted successfully"})
}

// UpdateShortURL implements URLService
func (h *URLHandler) UpdateShortUrl(ctx *gin.Context, shortPath string) {
	var req api.UpdateShortUrlJSONBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
		return
	}

	err := h.service.UpdateShortURL(ctx, req.OriginalUrl, ctx.Param("shortPath"), req.Expiry)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update short URL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Short URL updated successfully"})
}

// GetURLDetails implements URLService
func (h *URLHandler) GetShortUrlDetails(ctx *gin.Context, shortPath string) {
	urlDetails, err := h.service.GetURLDetails(ctx, shortPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get URL details"})
		return
	}
	if urlDetails == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, urlDetails)
}
func (h *URLHandler) GetShortUrlStats(ctx *gin.Context, shortPath string) {
	urlStats, err := h.urlStatService.GetURLStatistics(ctx, shortPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get URL statistics"})
		return
	}
	if urlStats == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, urlStats)
}
