package service

import (
	"context"

	"hadith-api-go/internal/domain/collection"
)

type collectionService struct {
	repo collection.CollectionRepository
}

func NewCollectionService(repo collection.CollectionRepository) collection.CollectionService {
	return &collectionService{repo: repo}
}

func (s *collectionService) List(ctx context.Context, limit, offset int) ([]collection.Collection, int, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *collectionService) FindByName(ctx context.Context, name string) (collection.Collection, error) {
	return s.repo.FindByName(ctx, name)
}
