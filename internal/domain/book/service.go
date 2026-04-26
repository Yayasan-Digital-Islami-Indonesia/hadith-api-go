package book

import "context"

// BookService defines the business operations for book data.
// Implement this interface in internal/service/book_service.go.
type BookService interface {
	List(ctx context.Context, collectionName string, limit, offset int) ([]Book, int, error)
	FindByNumber(ctx context.Context, collectionName, bookNumber string) (Book, error)
}
