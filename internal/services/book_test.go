package services

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/ydgi/hadith-api-go/internal/models"
	"github.com/ydgi/hadith-api-go/internal/repository"
	"gorm.io/gorm"
)

func TestBookService_GetAllBooks(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	db.AutoMigrate(&models.Book{})

	bookRepo := repository.NewBookRepository(db)
	service := NewBookService(bookRepo)

	bookRepo.Create(&models.Book{Slug: "bukhari", NameAr: "Test", NameEn: "Test", Totals: 10})

	books, err := service.GetAllBooks()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(books) != 1 {
		t.Errorf("expected 1 book, got %d", len(books))
	}
}

func TestBookService_GetBookBySlug(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	db.AutoMigrate(&models.Book{})

	bookRepo := repository.NewBookRepository(db)
	service := NewBookService(bookRepo)

	bookRepo.Create(&models.Book{Slug: "muslim", NameAr: "مسلم", NameEn: "Muslim", Totals: 20})

	book, err := service.GetBookBySlug("muslim")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if book == nil {
		t.Fatal("expected book, got nil")
	}
	if book.Slug != "muslim" {
		t.Errorf("expected slug 'muslim', got '%s'", book.Slug)
	}
}
