package services

import (
	"github.com/ydgi/hadith-api-go/internal/models"
	"github.com/ydgi/hadith-api-go/internal/repository"
)

type ChapterService struct {
	chapterRepo *repository.ChapterRepository
}

func NewChapterService(chapterRepo *repository.ChapterRepository) *ChapterService {
	return &ChapterService{chapterRepo: chapterRepo}
}

func (s *ChapterService) GetChaptersByBook(bookID uint) ([]models.Chapter, error) {
	return s.chapterRepo.GetByBookID(bookID)
}

func (s *ChapterService) GetChapter(id uint) (*models.Chapter, error) {
	return s.chapterRepo.GetByID(id)
}