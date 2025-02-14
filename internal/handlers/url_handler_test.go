package handlers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	api "url-shortener/generated"
	"url-shortener/internal/models"
	mocks "url-shortener/internal/services/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockExpiryTime, _ = time.Parse(time.RFC3339, "2024-05-30T12:34:56Z")

func setupHandler() (*mocks.URLService, *mocks.URLStatsService, *URLHandler) {
	mockURLService := mocks.URLService{}
	mockURLStatsService := mocks.URLStatsService{}
	return &mockURLService, &mockURLStatsService, NewURLHandler(&mockURLService, &mockURLStatsService)
}

func TestCreateShortURL_Success(t *testing.T) {
	mockURLService, mockURLStatsService, handler := setupHandler()
	mockURLService.On("CreateShortURL", mock.Anything, "https://www.example.com", &mockExpiryTime).Return("shortpath", nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	requestBody := api.CreateShortUrlJSONRequestBody{OriginalUrl: "https://www.example.com", Expiry: &mockExpiryTime}
	requestBodyBytes, _ := json.Marshal(requestBody)
	c.Request, _ = http.NewRequest(http.MethodPost, "/urls", bytes.NewBuffer(requestBodyBytes))

	handler.CreateShortUrl(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockURLService.AssertExpectations(t)
	mockURLStatsService.AssertExpectations(t)
}

func TestCreateShortURL_InvalidURL(t *testing.T) {
	_, _, handler := setupHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	requestBody := api.CreateShortUrlJSONRequestBody{OriginalUrl: "invalid_url"}
	requestBodyBytes, _ := json.Marshal(requestBody)
	c.Request, _ = http.NewRequest(http.MethodPost, "/urls", bytes.NewBuffer(requestBodyBytes))

	handler.CreateShortUrl(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRedirectToOriginalURL_Success(t *testing.T) {
	mockURLService, mockURLStatsService, handler := setupHandler()
	mockURLService.On("GetLongURL", mock.Anything, "shortpath").Return("https://www.example.com", nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/shortpath", nil)
	c.Params = gin.Params{{Key: "shortPath", Value: "shortpath"}}

	handler.RedirectToOriginalUrl(c, "shortpath")

	assert.Equal(t, http.StatusFound, w.Code)
	location, _ := url.Parse(w.Header().Get("Location"))
	assert.Equal(t, "https://www.example.com", location.String())
	mockURLService.AssertExpectations(t)
	mockURLStatsService.AssertExpectations(t)
}

func TestDeleteShortURL_Success(t *testing.T) {
	mockURLService, _, handler := setupHandler()
	mockURLService.On("DeleteURL", mock.Anything, "shortpath").Return(nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodDelete, "/shortpath", nil)
	c.Params = gin.Params{{Key: "shortPath", Value: "shortpath"}}

	handler.DeleteShortUrl(c, "shortpath")

	assert.Equal(t, http.StatusOK, w.Code)
	mockURLService.AssertExpectations(t)
}

func TestUpdateShortURL_Success(t *testing.T) {
	mockURLService, _, handler := setupHandler()
	mockURLService.On("UpdateShortURL", mock.Anything, "https://www.updated-example.com", "shortpath", &mockExpiryTime).Return(nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	requestBody := api.UpdateShortUrlJSONRequestBody{OriginalUrl: "https://www.updated-example.com", Expiry: &mockExpiryTime}
	requestBodyBytes, _ := json.Marshal(requestBody)
	c.Request, _ = http.NewRequest(http.MethodPut, "/shortpath", bytes.NewBuffer(requestBodyBytes))
	c.Params = gin.Params{{Key: "shortPath", Value: "shortpath"}}

	handler.UpdateShortUrl(c, "shortpath")

	assert.Equal(t, http.StatusOK, w.Code)
	mockURLService.AssertExpectations(t)
}

func TestGetShortURLStats_Success(t *testing.T) {
	_, mockURLStatsService, handler := setupHandler()
	mockStats := &models.URLStatistics{ShortPath: "shortpath", Last24Hours: 5, PastWeek: 5, AllTime: 5}
	mockURLStatsService.On("GetURLStatistics", mock.Anything, "shortpath").Return(mockStats, nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/shortpath/stats", nil)
	c.Params = gin.Params{{Key: "shortPath", Value: "shortpath"}}

	handler.GetShortUrlStats(c, "shortpath")

	assert.Equal(t, http.StatusOK, w.Code)
	mockURLStatsService.AssertExpectations(t)
}

func TestGetShortURLDetails_Success(t *testing.T) {
	mockURLService, _, handler := setupHandler()
	mockDetails := &models.URL{ShortPath: "shortpath", OriginalURL: "https://www.example.com", Expiry: &mockExpiryTime}
	mockURLService.On("GetURLDetails", mock.Anything, "shortpath").Return(mockDetails, nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/shortpath/details", nil)
	c.Params = gin.Params{{Key: "shortPath", Value: "shortpath"}}

	handler.GetShortUrlDetails(c, "shortpath")

	assert.Equal(t, http.StatusOK, w.Code)
	mockURLService.AssertExpectations(t)
}
func TestGetShortURLDetails_NotFound(t *testing.T) {
	mockURLService, _, handler := setupHandler()
	mockURLService.On("GetURLDetails", mock.Anything, "shortpath").Return(nil, nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/shortpath/details", nil)
	c.Params = gin.Params{{Key: "shortPath", Value: "shortpath"}}

	handler.GetShortUrlDetails(c, "shortpath")

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockURLService.AssertExpectations(t)
}

func TestGetShortURLStats_NotFound(t *testing.T) {
	_, mockURLStatsService, handler := setupHandler()
	mockURLStatsService.On("GetURLStatistics", mock.Anything, "shortpath").Return(nil, nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/shortpath/stats", nil)
	c.Params = gin.Params{{Key: "shortPath", Value: "shortpath"}}

	handler.GetShortUrlStats(c, "shortpath")

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockURLStatsService.AssertExpectations(t)
}

func TestRedirectToOriginalURL_NotFound(t *testing.T) {
	mockURLService, mockURLStatsService, handler := setupHandler()
	mockURLService.On("GetLongURL", mock.Anything, "shortpath").Return("", nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/shortpath", nil)
	c.Params = gin.Params{{Key: "shortPath", Value: "shortpath"}}

	handler.RedirectToOriginalUrl(c, "shortpath")

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockURLService.AssertExpectations(t)
	mockURLStatsService.AssertExpectations(t)
}
func TestDeleteShortURL_Failure(t *testing.T) {
	mockURLService, _, handler := setupHandler()
	mockURLService.On("DeleteURL", mock.Anything, "shortpath").Return(errors.New("failed to delete")).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodDelete, "/shortpath", nil)
	c.Params = gin.Params{{Key: "shortPath", Value: "shortpath"}}

	handler.DeleteShortUrl(c, "shortpath")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockURLService.AssertExpectations(t)
}

func TestUpdateShortURL_Failure(t *testing.T) {
	mockURLService, _, handler := setupHandler()
	mockURLService.On("UpdateShortURL", mock.Anything, "https://www.updated-example.com", "shortpath", &mockExpiryTime).Return(errors.New("failed to update")).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	requestBody := api.UpdateShortUrlJSONRequestBody{OriginalUrl: "https://www.updated-example.com", Expiry: &mockExpiryTime}
	requestBodyBytes, _ := json.Marshal(requestBody)
	c.Request, _ = http.NewRequest(http.MethodPut, "/shortpath", bytes.NewBuffer(requestBodyBytes))
	c.Params = gin.Params{{Key: "shortPath", Value: "shortpath"}}

	handler.UpdateShortUrl(c, "shortpath")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockURLService.AssertExpectations(t)
}

func TestGetShortURLStats_Failure(t *testing.T) {
	_, mockURLStatsService, handler := setupHandler()
	mockURLStatsService.On("GetURLStatistics", mock.Anything, "shortpath").Return(nil, errors.New("failed to get stats")).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/shortpath/stats", nil)
	c.Params = gin.Params{{Key: "shortPath", Value: "shortpath"}}

	handler.GetShortUrlStats(c, "shortpath")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockURLStatsService.AssertExpectations(t)
}

func TestGetShortURLDetails_Failure(t *testing.T) {
	mockURLService, _, handler := setupHandler()
	mockURLService.On("GetURLDetails", mock.Anything, "shortpath").Return(nil, errors.New("failed to get details")).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/shortpath/details", nil)
	c.Params = gin.Params{{Key: "shortPath", Value: "shortpath"}}

	handler.GetShortUrlDetails(c, "shortpath")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockURLService.AssertExpectations(t)
}
func TestUpdateShortURL_InvalidURL(t *testing.T) {
	_, _, handler := setupHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	requestBody := api.UpdateShortUrlJSONRequestBody{OriginalUrl: "invalid_url", Expiry: &mockExpiryTime}
	requestBodyBytes, _ := json.Marshal(requestBody)
	c.Request, _ = http.NewRequest(http.MethodPut, "/shortpath", bytes.NewBuffer(requestBodyBytes))
	c.Params = gin.Params{{Key: "shortPath", Value: "shortpath"}}

	handler.UpdateShortUrl(c, "shortpath")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateShortURL_InvalidRequest(t *testing.T) {
	_, _, handler := setupHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPut, "/shortpath", nil)
	c.Params = gin.Params{{Key: "shortPath", Value: "shortpath"}}

	handler.UpdateShortUrl(c, "shortpath")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestCreateShortURL_InvalidRequest(t *testing.T) {
	_, _, handler := setupHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/urls", nil)

	handler.CreateShortUrl(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRedirectToOriginalURL_Failure(t *testing.T) {
	mockURLService, mockURLStatsService, handler := setupHandler()
	mockURLService.On("GetLongURL", mock.Anything, "shortpath").Return("", errors.New("failed to get long URL")).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/shortpath", nil)
	c.Params = gin.Params{{Key: "shortPath", Value: "shortpath"}}

	handler.RedirectToOriginalUrl(c, "shortpath")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockURLService.AssertExpectations(t)
	mockURLStatsService.AssertExpectations(t)
}

func TestCreateShortURL_InternalError(t *testing.T) {
	mockURLService, mockURLStatsService, handler := setupHandler()
	mockURLService.On("CreateShortURL", mock.Anything, "https://www.example.com", &mockExpiryTime).Return("", errors.New("failed to create short URL")).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	requestBody := api.CreateShortUrlJSONRequestBody{OriginalUrl: "https://www.example.com", Expiry: &mockExpiryTime}
	requestBodyBytes, _ := json.Marshal(requestBody)
	c.Request, _ = http.NewRequest(http.MethodPost, "/urls", bytes.NewBuffer(requestBodyBytes))

	handler.CreateShortUrl(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockURLService.AssertExpectations(t)
	mockURLStatsService.AssertExpectations(t)
}

func TestCreateShortURL_SuccessHTTPS(t *testing.T) {
	mockURLService, mockURLStatsService, handler := setupHandler()
	mockURLService.On("CreateShortURL", mock.Anything, "https://www.example.com", &mockExpiryTime).Return("shortpath", nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	requestBody := api.CreateShortUrlJSONRequestBody{OriginalUrl: "https://www.example.com", Expiry: &mockExpiryTime}
	requestBodyBytes, _ := json.Marshal(requestBody)
	c.Request, _ = http.NewRequest(http.MethodPost, "/urls", bytes.NewBuffer(requestBodyBytes))
	c.Request.TLS = &tls.ConnectionState{}

	handler.CreateShortUrl(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockURLService.AssertExpectations(t)
	mockURLStatsService.AssertExpectations(t)
}

func TestValidateURL_InvalidHost(t *testing.T) {
	err := validateURL("http://")
	assert.Equal(t, ErrInvalidURLHost, err)
}
func TestValidateURL_InvalidScheme(t *testing.T) {
	err := validateURL("ftp://example.com")
	assert.Equal(t, ErrInvalidURLScheme, err)
}
