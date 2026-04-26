package hadith

import "context"

// HadithService defines the business operations for hadith data.
// Implement this interface in internal/service/hadith_service.go.
type HadithService interface {
	ListByBook(ctx context.Context, collectionName, bookNumber string, limit, offset int) ([]Hadith, int, error)
	FindByCollection(ctx context.Context, collectionName, hadithNumber string) (Hadith, error)
	Random(ctx context.Context) (Hadith, error)
}
