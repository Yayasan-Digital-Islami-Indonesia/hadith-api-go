package service

import (
	"context"

	"hadith-api-go/internal/domain/book"
)

type bookService struct {
	repo book.BookRepository
}

func NewBookService(repo book.BookRepository) book.BookService {
	return &bookService{repo: repo}
}

func (s *bookService) List(ctx context.Context, collectionName string, limit, offset int) ([]book.Book, int, error) {
	return s.repo.List(ctx, collectionName, limit, offset)
}

func (s *bookService) FindByNumber(ctx context.Context, collectionName, bookNumber string) (book.Book, error) {
	return s.repo.FindByNumber(ctx, collectionName, bookNumber)
}
