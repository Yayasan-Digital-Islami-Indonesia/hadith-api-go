package collection

import "context"

// CollectionService defines the business operations for hadith collection data.
// Implement this interface in internal/service/collection_service.go.
type CollectionService interface {
	List(ctx context.Context, limit, offset int) ([]Collection, int, error)
	FindByName(ctx context.Context, name string) (Collection, error)
}
