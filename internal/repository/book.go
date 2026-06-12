package repository

import (
	"fmt"

	"gorm.io/gorm"
	"github.com/ydgi/hadith-api-go/internal/models"
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