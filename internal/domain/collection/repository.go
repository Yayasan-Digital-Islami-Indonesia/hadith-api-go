package collection

import "context"

// CollectionRepository defines read-only access to hadith collection data.
// Implement this interface in internal/repository/collection_repository.go.
type CollectionRepository interface {
	List(ctx context.Context, limit, offset int) ([]Collection, int, error)
	FindByName(ctx context.Context, name string) (Collection, error)
}
