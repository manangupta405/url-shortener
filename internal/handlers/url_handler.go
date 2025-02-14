package handlers

import (
	"net/http"

	api "url-shortener/generated"
	"url-shortener/internal/repositories"
	"url-shortener/internal/services"
	"url-shortener/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type URLHandler struct {
	service        services.URLService
	urlStatService services.URLStatsService
	timeProvider   utils.TimeProvider
	api.ServerInterface
}

func NewURLHandler(service services.URLService, urlStatService services.URLStatsService, timeProvider utils.TimeProvider) *URLHandler {
	return &URLHandler{service: service, urlStatService: urlStatService, timeProvider: timeProvider}
}

func (h *URLHandler) CreateShortUrl(ctx *gin.Context) {
	var req api.CreateShortUrlJSONBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
		return
	}

	if err := validateURL(req.OriginalUrl); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	currentTime := h.timeProvider.Now()

	if req.Expiry.Before(currentTime) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Expiry date cannot be in the past"})
		return
	}
	shortPath, err := h.service.CreateShortURL(ctx, req.OriginalUrl, req.Expiry)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	shortUrl := scheme + "://" + ctx.Request.Host + "/" + shortPath

	response := &api.ShortenedUrlDetails{
		OriginalUrl: &req.OriginalUrl,
		ShortUrl:    &shortUrl,
	}

	ctx.JSON(http.StatusCreated, response)
}

func (h *URLHandler) RedirectToOriginalUrl(ctx *gin.Context, shortPath string) {
	longURL, err := h.service.GetLongURL(ctx, shortPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if longURL == "" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.Redirect(http.StatusFound, longURL)
}

// DeleteShortURL implements URLService
func (h *URLHandler) DeleteShortUrl(ctx *gin.Context, shortPath string) {
	err := h.service.DeleteURL(ctx, shortPath)
	if err == repositories.ErrShortURLNotFound {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Short URL not found"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
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
	if err := validateURL(req.OriginalUrl); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	currentTime := h.timeProvider.Now()

	if req.Expiry.Before(currentTime) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Expiry date cannot be in the past"})
		return
	}

	err := h.service.UpdateShortURL(ctx, req.OriginalUrl, shortPath, req.Expiry)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Short URL updated successfully"})
}

// GetURLDetails implements URLService
func (h *URLHandler) GetShortUrlDetails(ctx *gin.Context, shortPath string) {
	urlDetails, err := h.service.GetURLDetails(ctx, shortPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if urlStats == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, urlStats)
}
