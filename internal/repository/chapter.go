package repository

import (
	"fmt"

	"gorm.io/gorm"
	"github.com/ydgi/hadith-api-go/internal/models"
)

type ChapterRepository struct {
	db *gorm.DB
}

func NewChapterRepository(db *gorm.DB) *ChapterRepository {
	return &ChapterRepository{db: db}
}

func (r *ChapterRepository) GetByBookID(bookID uint) ([]models.Chapter, error) {
	var chapters []models.Chapter
	if err := r.db.Where("book_id = ?", bookID).Order("number ASC").Find(&chapters).Error; err != nil {
		return nil, fmt.Errorf("failed to get chapters: %w", err)
	}
	return chapters, nil
}

func (r *ChapterRepository) GetByID(id uint) (*models.Chapter, error) {
	var chapter models.Chapter
	if err := r.db.First(&chapter, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get chapter: %w", err)
	}
	return &chapter, nil
}

func (r *ChapterRepository) Create(chapter *models.Chapter) error {
	if err := r.db.Create(chapter).Error; err != nil {
		return fmt.Errorf("failed to create chapter: %w", err)
	}
	return nil
}