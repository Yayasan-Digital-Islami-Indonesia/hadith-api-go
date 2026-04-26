package hadith

import "context"

// HadithRepository defines read-only access to hadith data.
// Implement this interface in internal/repository/hadith_repository.go.
type HadithRepository interface {
	ListByBook(ctx context.Context, collectionName, bookNumber string, limit, offset int) ([]Hadith, int, error)
	FindByCollection(ctx context.Context, collectionName, hadithNumber string) (Hadith, error)
	Random(ctx context.Context) (Hadith, error)
}
