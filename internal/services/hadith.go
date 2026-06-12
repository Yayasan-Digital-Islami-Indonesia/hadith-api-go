package services

import (
	"github.com/ydgi/hadith-api-go/internal/models"
	"github.com/ydgi/hadith-api-go/internal/repository"
)

type HadithService struct {
	hadithRepo *repository.HadithRepository
}

func NewHadithService(hadithRepo *repository.HadithRepository) *HadithService {
	return &HadithService{hadithRepo: hadithRepo}
}

func (s *HadithService) GetHadith(id uint) (*models.Hadith, error) {
	return s.hadithRepo.GetByID(id)
}

func (s *HadithService) GetHadithByGlobalID(globalID string) (*models.Hadith, error) {
	return s.hadithRepo.GetByGlobalID(globalID)
}

func (s *HadithService) GetHadithByBookAndNumber(bookID uint, number int) (*models.Hadith, error) {
	return s.hadithRepo.GetByBookAndNumber(bookID, number)
}

func (s *HadithService) GetHadithsByChapter(chapterID uint, page, pageSize int) ([]models.Hadith, int64, error) {
	offset := (page - 1) * pageSize
	return s.hadithRepo.GetByChapterID(chapterID, pageSize, offset)
}

func (s *HadithService) GetRandomHadith() (*models.Hadith, error) {
	return s.hadithRepo.GetRandom()
}