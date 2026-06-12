package services

import (
	"strconv"

	"github.com/ydgi/hadith-api-go/internal/models"
	"github.com/ydgi/hadith-api-go/internal/repository"
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

func parseID(s string, id *uint) (bool, error) {
	parsed, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return false, err
	}
	*id = uint(parsed)
	return true, nil
}