package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/httputil"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r, // Здесь должен быть ваш обработчик запросов
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Создание канала для получения сигналов остановки
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		log.Println("Starting server...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Ожидание сигнала остановки
	<-sigChan

	// Создание контекста с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Остановка сервера с использованием graceful shutdown
	log.Println("Shutting down server...")

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")

}
