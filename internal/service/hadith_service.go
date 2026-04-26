package service

import (
	"context"

	"hadith-api-go/internal/domain/hadith"
)

type hadithService struct {
	repo hadith.HadithRepository
}

func NewHadithService(repo hadith.HadithRepository) hadith.HadithService {
	return &hadithService{repo: repo}
}

func (s *hadithService) ListByBook(ctx context.Context, collectionName, bookNumber string, limit, offset int) ([]hadith.Hadith, int, error) {
	return s.repo.ListByBook(ctx, collectionName, bookNumber, limit, offset)
}

func (s *hadithService) FindByCollection(ctx context.Context, collectionName, hadithNumber string) (hadith.Hadith, error) {
	return s.repo.FindByCollection(ctx, collectionName, hadithNumber)
}

func (s *hadithService) Random(ctx context.Context) (hadith.Hadith, error) {
	return s.repo.Random(ctx)
}
