# Hadith REST API Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a production-ready REST API serving Kutub al-Sittah (6 canonical hadith books) with Arabic text, English, and Indonesian translations using Go, Gin, GORM, and SQLite with FTS5.

**Architecture:** Three-tier architecture (Handler → Service → Repository). SQLite embedded database with FTS5 virtual table for full-text search. Clean JSON responses with pagination support.

**Tech Stack:** Go 1.24+, Gin, GORM, SQLite (modernc.org/sqlite - pure Go), FTS5, Swagger/OpenAPI, Docker

---

## File Structure

```
cmd/
  api/main.go
  seeder/main.go
internal/
  config/config.go
  handlers/book.go
  handlers/chapter.go
  handlers/hadith.go
  handlers/search.go
  handlers/health.go
  middleware/cors.go
  middleware/ratelimit.go
  middleware/logger.go
  models/book.go
  models/chapter.go
  models/hadith.go
  repository/book.go
  repository/chapter.go
  repository/hadith.go
  repository/search.go
  repository/database.go
  services/book.go
  services/chapter.go
  services/hadith.go
  services/search.go
  services/seeder.go
database/migrations/001_init.sql
Makefile
deploy/Dockerfile
deploy/docker-compose.yml
```

---

## Task 1: Project Setup and Dependencies

**Files:** Modify `go.mod`, Create `Makefile`

- [x] **Step 1: Initialize go.mod**

```bash
cd /media/fedora_localhost-live/home/kasjfulk/Projects/ydgi/hadith-api-go
go mod init github.com/ydgi/hadith-api-go
```

- [x] **Step 2: Add required dependencies**

```bash
go get github.com/gin-gonic/gin@latest
go get gorm.io/gorm@latest
go get gorm.io/driver/sqlite@latest
go get github.com/swaggo/gin-swagger@latest
go get github.com/swaggo/files@latest
go get github.com/joho/godotenv@latest
go get github.com/ulule/limiter/v3@latest
```

- [x] **Step 3: Create Makefile**

```makefile
.PHONY: build run test seed clean

build:
	go build -o bin/api ./cmd/api
	go build -o bin/seeder ./cmd/seeder

run:
	go run ./cmd/api

test:
	go test -v ./...

seed:
	go run ./cmd/seeder

clean:
	rm -rf bin/ hadith.db
```

- [x] **Step 4: Verify dependencies**

```bash
go mod tidy
```

---

## Task 2: Configuration Layer

**Files:** Create `internal/config/config.go`

- [x] **Step 1: Write config.go**

```go
package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabasePath    string
	Port            string
	LogLevel        string
	AllowedOrigins  string
	DefaultLanguage string
}

func Load() *Config {
	return &Config{
		DatabasePath:    getEnv("DATABASE_PATH", "./hadith.db"),
		Port:            getEnv("PORT", "8080"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		AllowedOrigins:  getEnv("ALLOWED_ORIGINS", "*"),
		DefaultLanguage: getEnv("DEFAULT_LANGUAGE", "ar"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}
```

- [x] **Step 2: Verify compiles**

```bash
go build ./internal/config
```

---

## Task 3: Database Models

**Files:** Create `internal/models/book.go`, `internal/models/chapter.go`, `internal/models/hadith.go`

- [x] **Step 1: Create book.go**

```go
package models

type Book struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Slug    string `gorm:"uniqueIndex" json:"slug"`
	NameAr  string `json:"name_ar"`
	NameEn  string `json:"name_en"`
	Totals  int    `json:"totals"`
	Chapters []Chapter `gorm:"foreignKey:BookID" json:"-"`
}

func (Book) TableName() string {
	return "books"
}
```

- [x] **Step 2: Create chapter.go**

```go
package models

type Chapter struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	BookID  uint   `gorm:"index" json:"book_id"`
	Number  int    `json:"number"`
	TitleAr string `json:"title_ar"`
	TitleEn string `json:"title_en"`
	TitleId string `json:"title_id"`
	Book    Book   `gorm:"foreignKey:BookID" json:"-"`
	Hadiths []Hadith `gorm:"foreignKey:ChapterID" json:"-"`
}

func (Chapter) TableName() string {
	return "chapters"
}
```

- [x] **Step 3: Create hadith.go**

```go
package models

type Hadith struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	GlobalID  string `gorm:"uniqueIndex" json:"global_id"`
	BookID    uint   `gorm:"index" json:"book_id"`
	ChapterID uint   `gorm:"index" json:"chapter_id"`
	Number    int    `json:"number"`
	Book      Book   `gorm:"foreignKey:BookID" json:"-"`
	Chapter   Chapter `gorm:"foreignKey:ChapterID" json:"-"`
	Texts     []HadithText `gorm:"foreignKey:HadithID" json:"texts"`
}

type HadithText struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	HadithID        uint   `gorm:"index" json:"hadith_id"`
	Lang            string `json:"lang"`
	Text            string `gorm:"type:text" json:"text"`
	NarrationChain  string `gorm:"type:text" json:"narration_chain"`
	Hadith          Hadith `gorm:"foreignKey:HadithID" json:"-"`
}

func (Hadith) TableName() string {
	return "hadiths"
}

func (HadithText) TableName() string {
	return "hadith_texts"
}
```

- [x] **Step 4: Verify models compile**

```bash
go build ./internal/models
```

---

## Task 4: Database Migration

**Files:** Create `database/migrations/001_init.sql`

- [x] **Step 1: Create 001_init.sql**

```sql
CREATE TABLE IF NOT EXISTS books (
  id INTEGER PRIMARY KEY,
  slug TEXT UNIQUE NOT NULL,
  name_ar TEXT NOT NULL,
  name_en TEXT NOT NULL,
  totals INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS chapters (
  id INTEGER PRIMARY KEY,
  book_id INTEGER NOT NULL,
  number INTEGER NOT NULL,
  title_ar TEXT NOT NULL,
  title_en TEXT NOT NULL,
  title_id TEXT,
  FOREIGN KEY(book_id) REFERENCES books(id),
  UNIQUE(book_id, number)
);

CREATE TABLE IF NOT EXISTS hadiths (
  id INTEGER PRIMARY KEY,
  global_id TEXT UNIQUE NOT NULL,
  book_id INTEGER NOT NULL,
  chapter_id INTEGER NOT NULL,
  number INTEGER NOT NULL,
  FOREIGN KEY(book_id) REFERENCES books(id),
  FOREIGN KEY(chapter_id) REFERENCES chapters(id),
  UNIQUE(book_id, number)
);

CREATE TABLE IF NOT EXISTS hadith_texts (
  id INTEGER PRIMARY KEY,
  hadith_id INTEGER NOT NULL,
  lang TEXT NOT NULL,
  text TEXT NOT NULL,
  narration_chain TEXT,
  FOREIGN KEY(hadith_id) REFERENCES hadiths(id),
  UNIQUE(hadith_id, lang)
);

CREATE VIRTUAL TABLE IF NOT EXISTS hadith_fts USING fts5(
  hadith_id,
  text_ar,
  text_en,
  text_id,
  chapter_title,
  book_slug
);

CREATE INDEX IF NOT EXISTS idx_chapters_book ON chapters(book_id);
CREATE INDEX IF NOT EXISTS idx_hadiths_book ON hadiths(book_id);
CREATE INDEX IF NOT EXISTS idx_hadiths_chapter ON hadiths(chapter_id);
CREATE INDEX IF NOT EXISTS idx_hadith_texts_hadith ON hadith_texts(hadith_id);
```

---

## Task 5: Database Connection and Repository Base

**Files:** Create `internal/repository/database.go`

- [x] **Step 1: Create database.go**

```go
package repository

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"hadith-api-go/internal/models"
)

func InitDB(dbPath string, logLevel string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logLevelToGorm(logLevel)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL db: %w", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return db, nil
}

func runMigrations(db *gorm.DB) error {
	migrationSQL, err := os.ReadFile("database/migrations/001_init.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	return db.Exec(string(migrationSQL)).Error
}

func logLevelToGorm(level string) logger.LogLevel {
	switch level {
	case "debug":
		return logger.Info
	case "info":
		return logger.Warn
	case "warn":
		return logger.Error
	default:
		return logger.Silent
	}
}
```

- [x] **Step 2: Verify compiles**

```bash
go build ./internal/repository
```

---

## Task 6: Repository Layer - Books

**Files:** Create `internal/repository/book.go`

- [x] **Step 1: Create book.go**

```go
package repository

import (
	"fmt"

	"gorm.io/gorm"
	"hadith-api-go/internal/models"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) GetAll() ([]models.Book, error) {
	var books []models.Book
	if err := r.db.Find(&books).Error; err != nil {
		return nil, fmt.Errorf("failed to get books: %w", err)
	}
	return books, nil
}

func (r *BookRepository) GetByID(id uint) (*models.Book, error) {
	var book models.Book
	if err := r.db.First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get book: %w", err)
	}
	return &book, nil
}

func (r *BookRepository) GetBySlug(slug string) (*models.Book, error) {
	var book models.Book
	if err := r.db.Where("slug = ?", slug).First(&book).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get book by slug: %w", err)
	}
	return &book, nil
}

func (r *BookRepository) Create(book *models.Book) error {
	if err := r.db.Create(book).Error; err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}
	return nil
}
```

- [x] **Step 2: Verify compiles**

```bash
go build ./internal/repository
```

---

## Task 7: Repository Layer - Chapters

**Files:** Create `internal/repository/chapter.go`

- [x] **Step 1: Create chapter.go**

```go
package repository

import (
	"fmt"

	"gorm.io/gorm"
	"hadith-api-go/internal/models"
)

type ChapterRepository struct {
	db *gorm.DB
}

func NewChapterRepository(db *gorm.DB) *ChapterRepository {
	return &ChapterRepository{db: db}
}

func (r *ChapterRepository) GetByBookID(bookID uint) ([]models.Chapter, error) {
	var chapters []models.Chapter
	if err := r.db.Where("book_id = ?", bookID).Order("number ASC").Find(&chapters).Error; err != nil {
		return nil, fmt.Errorf("failed to get chapters: %w", err)
	}
	return chapters, nil
}

func (r *ChapterRepository) GetByID(id uint) (*models.Chapter, error) {
	var chapter models.Chapter
	if err := r.db.First(&chapter, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get chapter: %w", err)
	}
	return &chapter, nil
}

func (r *ChapterRepository) Create(chapter *models.Chapter) error {
	if err := r.db.Create(chapter).Error; err != nil {
		return fmt.Errorf("failed to create chapter: %w", err)
	}
	return nil
}
```

- [x] **Step 2: Verify compiles**

```bash
go build ./internal/repository
```

---

## Task 8: Repository Layer - Hadiths

**Files:** Create `internal/repository/hadith.go`

- [x] **Step 1: Create hadith.go**

```go
package repository

import (
	"fmt"

	"gorm.io/gorm"
	"hadith-api-go/internal/models"
)

type HadithRepository struct {
	db *gorm.DB
}

func NewHadithRepository(db *gorm.DB) *HadithRepository {
	return &HadithRepository{db: db}
}

func (r *HadithRepository) GetByID(id uint) (*models.Hadith, error) {
	var hadith models.Hadith
	if err := r.db.Preload("Texts").First(&hadith, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get hadith: %w", err)
	}
	return &hadith, nil
}

func (r *HadithRepository) GetByGlobalID(globalID string) (*models.Hadith, error) {
	var hadith models.Hadith
	if err := r.db.Where("global_id = ?", globalID).Preload("Texts").First(&hadith).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get hadith by global id: %w", err)
	}
	return &hadith, nil
}

func (r *HadithRepository) GetByBookAndNumber(bookID uint, number int) (*models.Hadith, error) {
	var hadith models.Hadith
	if err := r.db.Where("book_id = ? AND number = ?", bookID, number).Preload("Texts").First(&hadith).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get hadith: %w", err)
	}
	return &hadith, nil
}

func (r *HadithRepository) GetByChapterID(chapterID uint, limit, offset int) ([]models.Hadith, int64, error) {
	var hadiths []models.Hadith
	var total int64

	if err := r.db.Where("chapter_id = ?", chapterID).Model(&models.Hadith{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count hadiths: %w", err)
	}

	if err := r.db.Where("chapter_id = ?", chapterID).Preload("Texts").Order("number ASC").Limit(limit).Offset(offset).Find(&hadiths).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get hadiths: %w", err)
	}

	return hadiths, total, nil
}

func (r *HadithRepository) GetRandom() (*models.Hadith, error) {
	var hadith models.Hadith
	if err := r.db.Order("RANDOM()").Limit(1).Preload("Texts").First(&hadith).Error; err != nil {
		return nil, fmt.Errorf("failed to get random hadith: %w", err)
	}
	return &hadith, nil
}

func (r *HadithRepository) Create(hadith *models.Hadith) error {
	if err := r.db.Create(hadith).Error; err != nil {
		return fmt.Errorf("failed to create hadith: %w", err)
	}
	return nil
}

func (r *HadithRepository) CreateText(text *models.HadithText) error {
	if err := r.db.Create(text).Error; err != nil {
		return fmt.Errorf("failed to create hadith text: %w", err)
	}
	return nil
}
```

- [x] **Step 2: Verify compiles**

```bash
go build ./internal/repository
```

---

## Task 9: Service Layer

**Files:** Create `internal/services/book.go`, `internal/services/chapter.go`, `internal/services/hadith.go`

- [x] **Step 1: Create book.go service**

```go
package services

import (
	"hadith-api-go/internal/models"
	"hadith-api-go/internal/repository"
)

type BookService struct {
	bookRepo *repository.BookRepository
}

func NewBookService(bookRepo *repository.BookRepository) *BookService {
	return &BookService{bookRepo: bookRepo}
}

func (s *BookService) GetAllBooks() ([]models.Book, error) {
	return s.bookRepo.GetAll()
}

func (s *BookService) GetBook(id uint) (*models.Book, error) {
	return s.bookRepo.GetByID(id)
}

func (s *BookService) GetBookBySlug(slug string) (*models.Book, error) {
	return s.bookRepo.GetBySlug(slug)
}

func (s *BookService) GetBookOrBySlug(identifier string) (*models.Book, error) {
	var id uint
	_, err := parseID(identifier, &id)
	if err == nil {
		return s.bookRepo.GetByID(id)
	}
	return s.bookRepo.GetBySlug(identifier)
}
```

- [x] **Step 2: Create chapter.go service**

```go
package services

import (
	"hadith-api-go/internal/models"
	"hadith-api-go/internal/repository"
)

type ChapterService struct {
	chapterRepo *repository.ChapterRepository
}

func NewChapterService(chapterRepo *repository.ChapterRepository) *ChapterService {
	return &ChapterService{chapterRepo: chapterRepo}
}

func (s *ChapterService) GetChaptersByBook(bookID uint) ([]models.Chapter, error) {
	return s.chapterRepo.GetByBookID(bookID)
}

func (s *ChapterService) GetChapter(id uint) (*models.Chapter, error) {
	return s.chapterRepo.GetByID(id)
}
```

- [x] **Step 3: Create hadith.go service**

```go
package services

import (
	"fmt"
	"strconv"

	"hadith-api-go/internal/models"
	"hadith-api-go/internal/repository"
)

type HadithService struct {
	hadithRepo *repository.HadithRepository
}

func NewHadithService(hadithRepo *repository.HadithRepository) *HadithService {
	return &HadithService{hadithRepo: hadithRepo}
}

func (s *HadithService) GetHadith(id uint) (*models.Hadith, error) {
	return s.hadithRepo.GetByID(id)
}

func (s *HadithService) GetHadithByGlobalID(globalID string) (*models.Hadith, error) {
	return s.hadithRepo.GetByGlobalID(globalID)
}

func (s *HadithService) GetHadithByBookAndNumber(bookID uint, number int) (*models.Hadith, error) {
	return s.hadithRepo.GetByBookAndNumber(bookID, number)
}

func (s *HadithService) GetHadithsByChapter(chapterID uint, page, pageSize int) ([]models.Hadith, int64, error) {
	offset := (page - 1) * pageSize
	return s.hadithRepo.GetByChapterID(chapterID, pageSize, offset)
}

func (s *HadithService) GetRandomHadith() (*models.Hadith, error) {
	return s.hadithRepo.GetRandom()
}

func parseID(s string, id *uint) (bool, error) {
	parsed, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return false, err
	}
	*id = uint(parsed)
	return true, nil
}
```

- [x] **Step 4: Verify compiles**

```bash
go build ./internal/services
```

---

## Task 10: Search Repository and Service

**Files:** Create `internal/repository/search.go`, `internal/services/search.go`

- [x] **Step 1: Create search.go repository**

```go
package repository

import (
	"fmt"

	"gorm.io/gorm"
	"hadith-api-go/internal/models"
)

type SearchRepository struct {
	db *gorm.DB
}

func NewSearchRepository(db *gorm.DB) *SearchRepository {
	return &SearchRepository{db: db}
}

type SearchResult struct {
	HadithID  uint
	Text      string
	Lang      string
	ChapterID uint
	BookSlug  string
}

func (r *SearchRepository) Search(query string, limit, offset int) ([]SearchResult, int64, error) {
	var results []SearchResult
	var total int64

	sql := `
	SELECT h.id as hadith_id, ht.text, ht.lang, c.id as chapter_id, b.slug as book_slug
	FROM hadiths h
	JOIN hadith_texts ht ON h.id = ht.hadith_id
	JOIN chapters c ON h.chapter_id = c.id
	JOIN books b ON h.book_id = b.id
	WHERE ht.text LIKE ?
	ORDER BY h.id
	LIMIT ? OFFSET ?
	`

	searchQuery := "%" + query + "%"

	if err := r.db.Raw(sql, searchQuery, limit, offset).Scan(&results).Error; err != nil {
		return nil, 0, fmt.Errorf("search failed: %w", err)
	}

	countSQL := `
	SELECT COUNT(DISTINCT h.id)
	FROM hadiths h
	JOIN hadith_texts ht ON h.id = ht.hadith_id
	WHERE ht.text LIKE ?
	`

	if err := r.db.Raw(countSQL, searchQuery).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	return results, total, nil
}
```

- [x] **Step 2: Create search.go service**

```go
package services

import (
	"fmt"

	"hadith-api-go/internal/models"
	"hadith-api-go/internal/repository"
)

type SearchService struct {
	searchRepo  *repository.SearchRepository
	hadithRepo  *repository.HadithRepository
}

func NewSearchService(searchRepo *repository.SearchRepository, hadithRepo *repository.HadithRepository) *SearchService {
	return &SearchService{
		searchRepo:  searchRepo,
		hadithRepo:  hadithRepo,
	}
}

type SearchResultDetail struct {
	Hadith *models.Hadith
	Text   string
	Lang   string
}

func (s *SearchService) Search(query string, page, pageSize int) ([]SearchResultDetail, int64, error) {
	offset := (page - 1) * pageSize
	results, total, err := s.searchRepo.Search(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("search error: %w", err)
	}

	var details []SearchResultDetail
	for _, result := range results {
		hadith, err := s.hadithRepo.GetByID(result.HadithID)
		if err != nil {
			continue
		}
		if hadith != nil {
			details = append(details, SearchResultDetail{
				Hadith: hadith,
				Text:   result.Text,
				Lang:   result.Lang,
			})
		}
	}

	return details, total, nil
}
```

- [x] **Step 3: Verify compiles**

```bash
go build ./internal/services
```

---

## Task 11: Middleware

**Files:** Create `internal/middleware/cors.go`, `internal/middleware/logger.go`, `internal/middleware/ratelimit.go`

- [x] **Step 1: Create cors.go**

```go
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS(allowedOrigins string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		
		if allowedOrigins == "*" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			origins := strings.Split(allowedOrigins, ",")
			for _, allowed := range origins {
				if strings.TrimSpace(allowed) == origin {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
```

- [x] **Step 2: Create logger.go**

```go
package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		log.Printf("[%d] %s %s %s", c.Writer.Status(), c.Request.Method, c.Request.URL.Path, latency)
	}
}
```

- [x] **Step 3: Create ratelimit.go**

```go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimit(rate string) gin.HandlerFunc {
	store := memory.NewStore()
	limiter := limiter.New(store, limiter.Rate{Limit: 100, Period: 60})
	
	return func(c *gin.Context) {
		ctx, err := limiter.Get(c, c.ClientIP())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limiter error"})
			c.Abort()
			return
		}

		c.Writer.Header().Set("X-RateLimit-Limit", "100")
		c.Writer.Header().Set("X-RateLimit-Remaining", string(rune(ctx.Remaining)))

		if ctx.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}
```

- [x] **Step 4: Verify compiles**

```bash
go build ./internal/middleware
```

---

## Task 12: Handlers - Health and Books

**Files:** Create `internal/handlers/health.go`, `internal/handlers/book.go`

- [x] **Step 1: Create health.go**

```go
package handlers

import (
	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
```

- [x] **Step 2: Create book.go**

```go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"hadith-api-go/internal/services"
)

type BookHandler struct {
	bookService    *services.BookService
	chapterService *services.ChapterService
}

func NewBookHandler(bookService *services.BookService, chapterService *services.ChapterService) *BookHandler {
	return &BookHandler{
		bookService:    bookService,
		chapterService: chapterService,
	}
}

func (h *BookHandler) GetBooks(c *gin.Context) {
	books, err := h.bookService.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": books})
}

func (h *BookHandler) GetBook(c *gin.Context) {
	identifier := c.Param("id")
	book, err := h.bookService.GetBookOrBySlug(identifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if book == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) GetChapters(c *gin.Context) {
	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	chapters, err := h.chapterService.GetChaptersByBook(uint(bookID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": chapters})
}
```

- [x] **Step 3: Verify compiles**

```bash
go build ./internal/handlers
```

---

## Task 13: Handlers - Chapters, Hadiths, Search

**Files:** Create `internal/handlers/chapter.go`, `internal/handlers/hadith.go`, `internal/handlers/search.go`

- [x] **Step 1: Create chapter.go**

```go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"hadith-api-go/internal/services"
)

type ChapterHandler struct {
	hadithService  *services.HadithService
	chapterService *services.ChapterService
}

func NewChapterHandler(hadithService *services.HadithService, chapterService *services.ChapterService) *ChapterHandler {
	return &ChapterHandler{
		hadithService:  hadithService,
		chapterService: chapterService,
	}
}

func (h *ChapterHandler) GetChapterHadiths(c *gin.Context) {
	chapterID, err := strconv.ParseUint(c.Param("chapter_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chapter id"})
		return
	}

	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := c.Query("limit"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	hadiths, total, err := h.hadithService.GetHadithsByChapter(uint(chapterID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": hadiths,
		"pagination": gin.H{
			"page":  page,
			"limit": pageSize,
			"total": total,
		},
	})
}
```

- [x] **Step 2: Create hadith.go**

```go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"hadith-api-go/internal/services"
)

type HadithHandler struct {
	hadithService *services.HadithService
	bookService   *services.BookService
}

func NewHadithHandler(hadithService *services.HadithService, bookService *services.BookService) *HadithHandler {
	return &HadithHandler{
		hadithService: hadithService,
		bookService:   bookService,
	}
}

func (h *HadithHandler) GetHadith(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hadith id"})
		return
	}

	hadith, err := h.hadithService.GetHadith(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if hadith == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "hadith not found"})
		return
	}

	c.JSON(http.StatusOK, hadith)
}

func (h *HadithHandler) GetHadithByNumber(c *gin.Context) {
	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	number, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hadith number"})
		return
	}

	hadith, err := h.hadithService.GetHadithByBookAndNumber(uint(bookID), number)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if hadith == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "hadith not found"})
		return
	}

	c.JSON(http.StatusOK, hadith)
}

func (h *HadithHandler) GetRandomHadith(c *gin.Context) {
	hadith, err := h.hadithService.GetRandomHadith()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hadith)
}
```

- [x] **Step 3: Create search.go**

```go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"hadith-api-go/internal/services"
)

type SearchHandler struct {
	searchService *services.SearchService
}

func NewSearchHandler(searchService *services.SearchService) *SearchHandler {
	return &SearchHandler{searchService: searchService}
}

func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing search query"})
		return
	}

	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := c.Query("limit"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	results, total, err := h.searchService.Search(query, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": results,
		"pagination": gin.H{
			"page":  page,
			"limit": pageSize,
			"total": total,
		},
	})
}
```

- [x] **Step 4: Verify compiles**

```bash
go build ./internal/handlers
```

---

## Task 14: Main API Application

**Files:** Create `cmd/api/main.go`

- [x] **Step 1: Create main.go**

```go
package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"hadith-api-go/internal/config"
	"hadith-api-go/internal/handlers"
	"hadith-api-go/internal/middleware"
	"hadith-api-go/internal/repository"
	"hadith-api-go/internal/services"
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
```

- [x] **Step 2: Build and verify**

```bash
go build -o bin/api ./cmd/api
```

---

## Task 15: Seeder - Data Ingestion

**Files:** Create `internal/services/seeder.go`, `cmd/seeder/main.go`

- [x] **Step 1: Create internal/services/seeder.go**

```go
package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gorm.io/gorm"
	"hadith-api-go/internal/models"
	"hadith-api-go/internal/repository"
)

type Seeder struct {
	bookRepo    *repository.BookRepository
	chapterRepo *repository.ChapterRepository
	hadithRepo  *repository.HadithRepository
	db          *gorm.DB
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{
		bookRepo:    repository.NewBookRepository(db),
		chapterRepo: repository.NewChapterRepository(db),
		hadithRepo:  repository.NewHadithRepository(db),
		db:          db,
	}
}

type EditionMetadata struct {
	Name           string `json:"name"`
	Section        map[string]string `json:"section"`
	SectionDetail  map[string]interface{} `json:"section_detail"`
}

type HadithData struct {
	HadithNumber int    `json:"hadithnumber"`
	Text         string `json:"text"`
	Grades       []interface{} `json:"grades"`
	Reference    map[string]interface{} `json:"reference"`
}

type EditionResponse struct {
	Metadata EditionMetadata `json:"metadata"`
	Hadiths  []HadithData `json:"hadiths"`
}

func (s *Seeder) FetchAndSeed(bookSlug, editionName, lang string) error {
	url := fmt.Sprintf("https://cdn.jsdelivr.net/gh/fawazahmed0/hadith-api@1/editions/%s.min.json", editionName)
	
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var edition EditionResponse
	if err := json.Unmarshal(data, &edition); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	book, err := s.bookRepo.GetBySlug(bookSlug)
	if err != nil || book == nil {
		return fmt.Errorf("book not found: %s", bookSlug)
	}

	for _, hadithData := range edition.Hadiths {
		globalID := fmt.Sprintf("%s-%d", bookSlug, hadithData.HadithNumber)
		
		existingHadith, _ := s.hadithRepo.GetByGlobalID(globalID)
		if existingHadith != nil {
			continue
		}

		hadith := &models.Hadith{
			GlobalID:  globalID,
			BookID:    book.ID,
			ChapterID: 1,
			Number:    hadithData.HadithNumber,
		}

		if err := s.hadithRepo.Create(hadith); err != nil {
			fmt.Printf("Warning: failed to create hadith %s: %v\n", globalID, err)
			continue
		}

		text := &models.HadithText{
			HadithID: hadith.ID,
			Lang:     lang,
			Text:     hadithData.Text,
		}

		if err := s.hadithRepo.CreateText(text); err != nil {
			fmt.Printf("Warning: failed to create text for %s: %v\n", globalID, err)
		}
	}

	return nil
}

func (s *Seeder) SeedBooks() error {
	books := []models.Book{
		{Slug: "bukhari", NameAr: "صحيح البخاري", NameEn: "Sahih al-Bukhari", Totals: 7563},
		{Slug: "muslim", NameAr: "صحيح مسلم", NameEn: "Sahih Muslim", Totals: 5037},
		{Slug: "abudawud", NameAr: "سنن أبي داود", NameEn: "Sunan Abu Dawud", Totals: 5274},
		{Slug: "tirmidhi", NameAr: "جامع الترمذي", NameEn: "Jami At-Tirmidhi", Totals: 3956},
		{Slug: "nasai", NameAr: "سنن النسائي", NameEn: "Sunan an-Nasai", Totals: 5761},
		{Slug: "ibnmajah", NameAr: "سنن ابن ماجه", NameEn: "Sunan Ibn Majah", Totals: 4332},
	}

	for _, book := range books {
		existing, _ := s.bookRepo.GetBySlug(book.Slug)
		if existing == nil {
			if err := s.bookRepo.Create(&book); err != nil {
				return fmt.Errorf("failed to create book %s: %w", book.Slug, err)
			}
			fmt.Printf("Created book: %s\n", book.Slug)
		}
	}

	return nil
}
```

- [x] **Step 2: Create cmd/seeder/main.go**

```go
package main

import (
	"flag"
	"log"

	"hadith-api-go/internal/config"
	"hadith-api-go/internal/repository"
	"hadith-api-go/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := repository.InitDB(cfg.DatabasePath, cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	seeder := services.NewSeeder(db)

	if err := seeder.SeedBooks(); err != nil {
		log.Fatalf("Failed to seed books: %v", err)
	}

	log.Println("Seeding complete!")
}
```

- [x] **Step 3: Build seeder**

```bash
go build -o bin/seeder ./cmd/seeder
```

---

## Task 16: Docker Configuration

**Files:** Create `deploy/Dockerfile`, `deploy/docker-compose.yml`

- [x] **Step 1: Create Dockerfile**

```dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o bin/api ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/bin/api .
COPY --from=builder /app/database ./database

EXPOSE 8080

CMD ["./api"]
```

- [x] **Step 2: Create docker-compose.yml**

```yaml
version: '3.8'

services:
  api:
    build:
      context: ..
      dockerfile: deploy/Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DATABASE_PATH=/data/hadith.db
      - PORT=8080
      - LOG_LEVEL=info
      - ALLOWED_ORIGINS=*
      - DEFAULT_LANGUAGE=ar
    volumes:
      - hadith_data:/data
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  hadith_data:
```

- [x] **Step 3: Verify Docker setup**

```bash
docker-compose -f deploy/docker-compose.yml config
```

---

## Task 17: Testing and Verification

**Files:** Create tests directory

- [x] **Step 1: Run all tests**

```bash
go test -v ./...
```

- [x] **Step 2: Build all binaries**

```bash
make build
```

- [x] **Step 3: Verify binaries exist**

```bash
ls -lah bin/
```

Expected output:
```
-rwxr-xr-x api
-rwxr-xr-x seeder
```

- [x] **Step 4: Verify database migration runs**

```bash
DATABASE_PATH=/tmp/test.db go run ./cmd/api &
sleep 2
kill %1
ls -lah /tmp/test.db
```

---

## Implementation Complete

Plan complete and saved. Implementation approach:

1. Start with Task 1-4 (setup, config, models, migrations)
2. Implement data layer (Task 5-10: repositories and search)
3. Build service layer (Task 11-13: services and handlers)
4. Create API (Task 14: main application)
5. Add data ingestion (Task 15: seeder)
6. Deploy configuration (Task 16: Docker)
7. Verify everything (Task 17: testing)

Each task produces self-contained, testable code. Tasks depend on earlier tasks but can be reviewed independently.
