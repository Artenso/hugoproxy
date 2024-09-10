package main

import (
	"net/http"
	_ "net/http/httputil"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware

	_ "test/docs"
)

const (
	apiKey    = ""
	secretKey = ""
)

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(Secret), nil)
}

// @title Proxy
// @version 1.0
// @description Documentation of Proxy API.

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	geoserv := NewGeoService(apiKey, secretKey)

	storage := NewStorage()

	handler := NewHandler(geoserv, storage)

	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	r.Post("/api/login", handler.LoginHandler)
	r.Post("/api/register", handler.RegisterHandler)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/api/address/search", handler.AddressSearchHandler)
		r.Post("/api/address/geocode", handler.GeoCodeHandler)
	})

	http.ListenAndServe(":8080", r)
}
