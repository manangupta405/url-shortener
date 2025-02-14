// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Message *string `json:"message,omitempty"`
}

// ShortenedUrlDetails defines model for ShortenedUrlDetails.
type ShortenedUrlDetails struct {
	OriginalUrl *string `json:"originalUrl,omitempty"`
	ShortPath   *string `json:"short-path,omitempty"`
	ShortUrl    *string `json:"shortUrl,omitempty"`
}

// CreateShortUrlJSONBody defines parameters for CreateShortUrl.
type CreateShortUrlJSONBody struct {
	OriginalUrl string `json:"originalUrl"`
}

// UpdateShortUrlJSONBody defines parameters for UpdateShortUrl.
type UpdateShortUrlJSONBody struct {
	OriginalUrl string `json:"originalUrl"`
}

// CreateShortUrlJSONRequestBody defines body for CreateShortUrl for application/json ContentType.
type CreateShortUrlJSONRequestBody CreateShortUrlJSONBody

// UpdateShortUrlJSONRequestBody defines body for UpdateShortUrl for application/json ContentType.
type UpdateShortUrlJSONRequestBody UpdateShortUrlJSONBody
