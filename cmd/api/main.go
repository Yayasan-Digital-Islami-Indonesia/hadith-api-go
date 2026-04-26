package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"hadith-api-go/internal/config"
	"hadith-api-go/internal/database"
	"hadith-api-go/internal/handler"
	"hadith-api-go/internal/middleware"
	"hadith-api-go/internal/repository"
	"hadith-api-go/internal/service"
)

func main() {
	cfg := config.Load()
	setupLogger(cfg.LogLevel)

	db, err := database.New(cfg.DBPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close database")
		}
	}()

	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(middleware.Logging())
	if cfg.AllowedOrigins != "" {
		r.Use(middleware.CORS(cfg.AllowedOrigins))
	}

	// Health check
	healthCheckRepo := repository.NewHealthCheckRepository(db)
	healthCheckService := service.NewHealthCheckService(healthCheckRepo)
	healthCheckHandler := handler.NewHealthCheckHandler(healthCheckService)

	// Collections
	collectionRepo := repository.NewCollectionRepository(db)
	collectionService := service.NewCollectionService(collectionRepo)
	collectionHandler := handler.NewCollectionHandler(collectionService)

	// Books
	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookService)

	// Hadiths
	hadithRepo := repository.NewHadithRepository(db)
	hadithService := service.NewHadithService(hadithRepo)
	hadithHandler := handler.NewHadithHandler(hadithService)

	// Docs
	docsHandler := handler.NewDocsHandler()

	// Routes
	r.GET("/health", healthCheckHandler.HealthCheck)
	r.GET("/health/ready", healthCheckHandler.ReadyCheck)

	v1 := r.Group("/v1")
	{
		v1.GET("/collections", collectionHandler.List)
		v1.GET("/collections/:name", collectionHandler.Detail)
		v1.GET("/collections/:name/books", bookHandler.List)
		v1.GET("/collections/:name/books/:bookNumber", bookHandler.Detail)
		v1.GET("/collections/:name/books/:bookNumber/hadiths", hadithHandler.ListByBook)
		v1.GET("/collections/:name/hadiths/:hadithNumber", hadithHandler.FindByCollection)
		v1.GET("/hadiths/random", hadithHandler.Random)
	}

	// Documentation
	r.GET("/docs", docsHandler.ServeDocs)
	r.GET("/openapi.yaml", docsHandler.ServeOpenAPI)

	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Info().Str("addr", addr).Msg("starting server")
	if err := r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("server stopped")
	}
}

func setupLogger(level string) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(lvl)
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}
