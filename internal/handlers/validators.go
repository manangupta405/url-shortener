package handlers

import (
	"errors"
	"net/url"
)

var (
	ErrEmptyURL         = errors.New("URL cannot be empty")
	ErrInvalidURLFormat = errors.New("invalid URL format")
	ErrInvalidURLScheme = errors.New("invalid URL scheme (must be http or https)")
	ErrInvalidURLHost   = errors.New("invalid URL host")
)

func validateURL(urlStr string) error {
	if urlStr == "" {
		return ErrEmptyURL
	}

	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return ErrInvalidURLFormat
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return ErrInvalidURLScheme
	}

	if u.Host == "" {
		return ErrInvalidURLHost
	}

	return nil
}
