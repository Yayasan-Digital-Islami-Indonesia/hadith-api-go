package repository

import (
	"fmt"

	"gorm.io/gorm"
	"github.com/ydgi/hadith-api-go/internal/models"
)

type HadithRepository struct {
	db *gorm.DB
}

func NewHadithRepository(db *gorm.DB) *HadithRepository {
	return &HadithRepository{db: db}
}

func (r *HadithRepository) GetByID(id uint) (*models.Hadith, error) {
	var hadith models.Hadith
	if err := r.db.Preload("Texts").First(&hadith, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get hadith: %w", err)
	}
	return &hadith, nil
}

func (r *HadithRepository) GetByGlobalID(globalID string) (*models.Hadith, error) {
	var hadith models.Hadith
	if err := r.db.Where("global_id = ?", globalID).Preload("Texts").First(&hadith).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get hadith by global id: %w", err)
	}
	return &hadith, nil
}

func (r *HadithRepository) GetByBookAndNumber(bookID uint, number int) (*models.Hadith, error) {
	var hadith models.Hadith
	if err := r.db.Where("book_id = ? AND number = ?", bookID, number).Preload("Texts").First(&hadith).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get hadith: %w", err)
	}
	return &hadith, nil
}

func (r *HadithRepository) GetByChapterID(chapterID uint, limit, offset int) ([]models.Hadith, int64, error) {
	var hadiths []models.Hadith
	var total int64

	if err := r.db.Where("chapter_id = ?", chapterID).Model(&models.Hadith{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count hadiths: %w", err)
	}

	if err := r.db.Where("chapter_id = ?", chapterID).Preload("Texts").Order("number ASC").Limit(limit).Offset(offset).Find(&hadiths).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get hadiths: %w", err)
	}

	return hadiths, total, nil
}

func (r *HadithRepository) GetRandom() (*models.Hadith, error) {
	var hadith models.Hadith
	if err := r.db.Order("RANDOM()").Limit(1).Preload("Texts").First(&hadith).Error; err != nil {
		return nil, fmt.Errorf("failed to get random hadith: %w", err)
	}
	return &hadith, nil
}

func (r *HadithRepository) Create(hadith *models.Hadith) error {
	if err := r.db.Create(hadith).Error; err != nil {
		return fmt.Errorf("failed to create hadith: %w", err)
	}
	return nil
}

func (r *HadithRepository) CreateText(text *models.HadithText) error {
	if err := r.db.Create(text).Error; err != nil {
		return fmt.Errorf("failed to create hadith text: %w", err)
	}
	return nil
}