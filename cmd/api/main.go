package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ydgi/hadith-api-go/internal/config"
	"github.com/ydgi/hadith-api-go/internal/handlers"
	"github.com/ydgi/hadith-api-go/internal/middleware"
	"github.com/ydgi/hadith-api-go/internal/repository"
	"github.com/ydgi/hadith-api-go/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := repository.InitDB(cfg.DatabasePath, cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	bookRepo := repository.NewBookRepository(db)
	chapterRepo := repository.NewChapterRepository(db)
	hadithRepo := repository.NewHadithRepository(db)
	searchRepo := repository.NewSearchRepository(db)

	bookService := services.NewBookService(bookRepo)
	chapterService := services.NewChapterService(chapterRepo)
	hadithService := services.NewHadithService(hadithRepo)
	searchService := services.NewSearchService(searchRepo, hadithRepo)

	bookHandler := handlers.NewBookHandler(bookService, chapterService)
	chapterHandler := handlers.NewChapterHandler(hadithService, chapterService)
	hadithHandler := handlers.NewHadithHandler(hadithService, bookService)
	searchHandler := handlers.NewSearchHandler(searchService)

	router := gin.Default()

	router.Use(middleware.CORS(cfg.AllowedOrigins))
	router.Use(middleware.Logger())
	router.Use(middleware.RateLimit("100"))

	router.GET("/health", handlers.HealthCheck)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/books", bookHandler.GetBooks)
		v1.GET("/books/:id", bookHandler.GetBook)
		v1.GET("/books/:id/chapters", bookHandler.GetChapters)
		v1.GET("/books/:id/chapters/:chapter_id", chapterHandler.GetChapterHadiths)
		v1.GET("/hadith/:id", hadithHandler.GetHadith)
		v1.GET("/books/:id/hadith/:number", hadithHandler.GetHadithByNumber)
		v1.GET("/search", searchHandler.Search)
		v1.GET("/random", hadithHandler.GetRandomHadith)
	}

	addr := ":" + cfg.Port
	log.Printf("Starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}