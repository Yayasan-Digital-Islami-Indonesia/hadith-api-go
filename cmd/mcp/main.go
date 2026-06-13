package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/ydgi/hadith-api-go/internal/config"
	"github.com/ydgi/hadith-api-go/internal/repository"
	"github.com/ydgi/hadith-api-go/internal/services"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB
	bookService    *services.BookService
	chapterService *services.ChapterService
	hadithService  *services.HadithService
	searchService  *services.SearchService
)

func initServices() error {
	cfg := config.Load()
	var err error
	db, err = repository.InitDB(cfg.DatabasePath, cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to init database: %w", err)
	}

	bookRepo := repository.NewBookRepository(db)
	chapterRepo := repository.NewChapterRepository(db)
	hadithRepo := repository.NewHadithRepository(db)
	searchRepo := repository.NewSearchRepository(db)

	bookService = services.NewBookService(bookRepo)
	chapterService = services.NewChapterService(chapterRepo)
	hadithService = services.NewHadithService(hadithRepo)
	searchService = services.NewSearchService(searchRepo, hadithRepo)

	return nil
}

func main() {
	if err := initServices(); err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}

	srv := server.NewMCPServer(
		"hadith-api",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
	)

	// Tools
	srv.AddTool(mcp.NewTool("get_books",
		mcp.WithDescription("Get all hadith books"),
	), handleGetBooks)

	srv.AddTool(mcp.NewTool("get_book",
		mcp.WithDescription("Get a book by ID or slug"),
		mcp.WithString("identifier", mcp.Required(), mcp.Description("Book ID or slug")),
	), handleGetBook)

	srv.AddTool(mcp.NewTool("get_chapters",
		mcp.WithDescription("Get chapters for a book"),
		mcp.WithString("book_id", mcp.Required(), mcp.Description("Book ID")),
	), handleGetChapters)

	srv.AddTool(mcp.NewTool("get_chapter_hadiths",
		mcp.WithDescription("Get hadiths from a chapter"),
		mcp.WithNumber("chapter_id", mcp.Required(), mcp.Description("Chapter ID")),
		mcp.WithNumber("page", mcp.Description("Page number (default: 1)")),
		mcp.WithNumber("limit", mcp.Description("Items per page (default: 20, max: 100)")),
	), handleGetChapterHadiths)

	srv.AddTool(mcp.NewTool("get_hadith",
		mcp.WithDescription("Get a hadith by ID"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Hadith ID")),
	), handleGetHadith)

	srv.AddTool(mcp.NewTool("get_hadith_by_number",
		mcp.WithDescription("Get hadith by book and number"),
		mcp.WithNumber("book_id", mcp.Required(), mcp.Description("Book ID")),
		mcp.WithNumber("number", mcp.Required(), mcp.Description("Hadith number in book")),
	), handleGetHadithByNumber)

	srv.AddTool(mcp.NewTool("search_hadith",
		mcp.WithDescription("Search hadiths by text"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
		mcp.WithNumber("page", mcp.Description("Page number (default: 1)")),
		mcp.WithNumber("limit", mcp.Description("Items per page (default: 20, max: 100)")),
	), handleSearch)

	srv.AddTool(mcp.NewTool("get_random_hadith",
		mcp.WithDescription("Get a random hadith"),
	), handleGetRandomHadith)

	if err := server.ServeStdio(srv); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
