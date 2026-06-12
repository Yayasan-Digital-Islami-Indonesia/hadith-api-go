package services

import (
	"fmt"

	"github.com/ydgi/hadith-api-go/internal/models"
	"github.com/ydgi/hadith-api-go/internal/repository"
)

type SearchService struct {
	searchRepo *repository.SearchRepository
	hadithRepo *repository.HadithRepository
}

func NewSearchService(searchRepo *repository.SearchRepository, hadithRepo *repository.HadithRepository) *SearchService {
	return &SearchService{
		searchRepo: searchRepo,
		hadithRepo: hadithRepo,
	}
}

type SearchResultDetail struct {
	Hadith *models.Hadith
	Text   string
	Lang   string
}

func (s *SearchService) Search(query string, page, pageSize int) ([]SearchResultDetail, int64, error) {
	offset := (page - 1) * pageSize
	results, total, err := s.searchRepo.Search(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("search error: %w", err)
	}

	var details []SearchResultDetail
	for _, result := range results {
		hadith, err := s.hadithRepo.GetByID(result.HadithID)
		if err != nil {
			continue
		}
		if hadith != nil {
			details = append(details, SearchResultDetail{
				Hadith: hadith,
				Text:   result.Text,
				Lang:   result.Lang,
			})
		}
	}

	return details, total, nil
}