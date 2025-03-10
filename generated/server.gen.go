// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create a shortened URL
	// (POST /urls)
	CreateShortUrl(c *gin.Context)
	// Delete a shortened URL
	// (DELETE /urls/{short-path})
	DeleteShortUrl(c *gin.Context, shortPath string)
	// Retrieve details of a shortened URL
	// (GET /urls/{short-path})
	GetShortUrlDetails(c *gin.Context, shortPath string)
	// Update a shortened URL
	// (PUT /urls/{short-path})
	UpdateShortUrl(c *gin.Context, shortPath string)
	// Get access statistics for a shortened URL
	// (GET /urls/{short-path}/stats)
	GetShortUrlStats(c *gin.Context, shortPath string)
	// Redirect to the original URL
	// (GET /{short-path})
	RedirectToOriginalUrl(c *gin.Context, shortPath string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// CreateShortUrl operation middleware
func (siw *ServerInterfaceWrapper) CreateShortUrl(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.CreateShortUrl(c)
}

// DeleteShortUrl operation middleware
func (siw *ServerInterfaceWrapper) DeleteShortUrl(c *gin.Context) {

	var err error

	// ------------- Path parameter "short-path" -------------
	var shortPath string

	err = runtime.BindStyledParameterWithOptions("simple", "short-path", c.Param("short-path"), &shortPath, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter short-path: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.DeleteShortUrl(c, shortPath)
}

// GetShortUrlDetails operation middleware
func (siw *ServerInterfaceWrapper) GetShortUrlDetails(c *gin.Context) {

	var err error

	// ------------- Path parameter "short-path" -------------
	var shortPath string

	err = runtime.BindStyledParameterWithOptions("simple", "short-path", c.Param("short-path"), &shortPath, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter short-path: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetShortUrlDetails(c, shortPath)
}

// UpdateShortUrl operation middleware
func (siw *ServerInterfaceWrapper) UpdateShortUrl(c *gin.Context) {

	var err error

	// ------------- Path parameter "short-path" -------------
	var shortPath string

	err = runtime.BindStyledParameterWithOptions("simple", "short-path", c.Param("short-path"), &shortPath, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter short-path: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.UpdateShortUrl(c, shortPath)
}

// GetShortUrlStats operation middleware
func (siw *ServerInterfaceWrapper) GetShortUrlStats(c *gin.Context) {

	var err error

	// ------------- Path parameter "short-path" -------------
	var shortPath string

	err = runtime.BindStyledParameterWithOptions("simple", "short-path", c.Param("short-path"), &shortPath, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter short-path: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetShortUrlStats(c, shortPath)
}

// RedirectToOriginalUrl operation middleware
func (siw *ServerInterfaceWrapper) RedirectToOriginalUrl(c *gin.Context) {

	var err error

	// ------------- Path parameter "short-path" -------------
	var shortPath string

	err = runtime.BindStyledParameterWithOptions("simple", "short-path", c.Param("short-path"), &shortPath, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter short-path: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.RedirectToOriginalUrl(c, shortPath)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.POST(options.BaseURL+"/urls", wrapper.CreateShortUrl)
	router.DELETE(options.BaseURL+"/urls/:short-path", wrapper.DeleteShortUrl)
	router.GET(options.BaseURL+"/urls/:short-path", wrapper.GetShortUrlDetails)
	router.PUT(options.BaseURL+"/urls/:short-path", wrapper.UpdateShortUrl)
	router.GET(options.BaseURL+"/urls/:short-path/stats", wrapper.GetShortUrlStats)
	router.GET(options.BaseURL+"/:short-path", wrapper.RedirectToOriginalUrl)
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xYS4/bNhD+KwTbozZ2HBdIddsmQbrA9gE/0EOwh4k0sphSJEuOtjEM//diKMkvKSnc",
	"BN4suqf1jjTkx5nvmxlxIzNbOWvQUJDpRoasxArizzfeWz/D4KwJyAbnrUNPCuPjCkOAVXxAa4cylYG8",
	"Miu53Sadxb7/gBnJbSLnpfWEBvOl16+RQOnQXxI/OuXX/KuwvgKSqcyB8IpUhTI53SaR1quVMqCXXg/A",
	"SGTgTa8cUPnpx8O+Q0dYzm7nBKQCqWwAPGi9YJzpRuYYMq8cKWtkKheWQAtTV+/RC1sIyDIMAcP+RMoQ",
	"rtDzJhoCTaY/29qH/kq/9tYQyggqUbCbmExFGR2HFnYQ6A/EP89ZlX3E3+zUX7EfITYpU9j+Dtdi9ma+",
	"ENe/34jCehEaLiizEsvZbUhEBQZW/C+VWCUCTC5KMLlmk8dcecwoPGMQijRvuZzdio5RnteVibxHH5rt",
	"nj8bPxtHfjg04JRM5Yto4ihQGQM7qn3LQBuI/3IqgfHe5DKVrzwC4bwjSCI9/lVjoJ9sHvmZWUNooiM4",
	"p1UWXUcfAgPoVPTVCb5zqr3qv75tUCqPuUzfHfneDSZr/zb5GqOhEXvEOhk//4KTnqUtNh0zZlcumCEi",
	"i9nIOSLT8fgsVN97LGQqvxvty9yorXGj4wI3gOLG3INWuWiTLxystYUWx4+Xw/HKmkKrjMRVDAdoj5Cv",
	"Ox21gZlcENDCWhbtuotMEFdiBoRCq0qRwI8ZYt7g+uGyCSP0BrQI6O/RC2SHyLhQVxWw9FppC9iHj4PK",
	"coJVYOFwiH/hioQVI75j91guRpt9Q9k2VU4jYb96vI72g+rhwEOFhFzU322kYqyxLSXSQBXFsW9Vp8JM",
	"DqJzKqW7nmin/fLLB2qg5iLUscYXtdbrhsfTy6WHgRhLorC1eSLtOaRtGHUOaRO5woHG9hap42U3hl2U",
	"nl8vrkMT5Sc4lzePhUfyCu+7TvLE/EfA/Fmbs10SeU49RwauHpDB0uVwwQr9/5sdL67z44mxjvn9xibG",
	"B603DzivNskQWWv9JstMUxC+eCocBYLmIuXfuu88vvhIe+/xVchAyK/jlCnC7qUH7L3HleGpC/8XebxF",
	"am+HDnNaWP8ZxRwwpFHL6efToERm7X3Pwv520IQuqJMX40n/E6pDFS+prOj6Y3viEiHH5sLu1jZJGv4K",
	"I7u7zxJk5edgbZ9G1Mczou5SGq9NT9jR6aF7jdXA/nHBhs01z2iyJHLpaKRtBrq0gdKX45djub3b/hMA",
	"AP//sQScXaUXAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
