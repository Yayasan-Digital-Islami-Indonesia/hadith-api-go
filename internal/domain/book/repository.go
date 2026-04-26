package book

import "context"

// BookRepository defines read-only access to book data within a collection.
// Implement this interface in internal/repository/book_repository.go.
type BookRepository interface {
	List(ctx context.Context, collectionName string, limit, offset int) ([]Book, int, error)
	FindByNumber(ctx context.Context, collectionName, bookNumber string) (Book, error)
}
