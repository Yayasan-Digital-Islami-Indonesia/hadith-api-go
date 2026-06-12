package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"gorm.io/gorm"
	"github.com/ydgi/hadith-api-go/internal/models"
	"github.com/ydgi/hadith-api-go/internal/repository"
)

type Seeder struct {
	bookRepo    *repository.BookRepository
	chapterRepo *repository.ChapterRepository
	hadithRepo  *repository.HadithRepository
	db          *gorm.DB
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{
		bookRepo:    repository.NewBookRepository(db),
		chapterRepo: repository.NewChapterRepository(db),
		hadithRepo:  repository.NewHadithRepository(db),
		db:          db,
	}
}

type EditionMetadata struct {
	Name          string                 `json:"name"`
	Section       map[string]string      `json:"section"`
	SectionDetail  map[string]interface{} `json:"section_detail"`
}

type HadithData struct {
	HadithNumber interface{}            `json:"hadithnumber"`
	Text         string                 `json:"text"`
	Grades       []interface{}          `json:"grades"`
	Reference    map[string]interface{} `json:"reference"`
}

type EditionResponse struct {
	Metadata EditionMetadata `json:"metadata"`
	Hadiths  []HadithData    `json:"hadiths"`
}

func (s *Seeder) FetchAndSeed(bookSlug, editionName, lang string) error {
	url := fmt.Sprintf("https://cdn.jsdelivr.net/gh/fawazahmed0/hadith-api@1/editions/%s.min.json", editionName)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var edition EditionResponse
	if err := json.Unmarshal(data, &edition); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	book, err := s.bookRepo.GetBySlug(bookSlug)
	if err != nil || book == nil {
		return fmt.Errorf("book not found: %s", bookSlug)
	}

	for _, hadithData := range edition.Hadiths {
		var hadithNum int
		switch v := hadithData.HadithNumber.(type) {
		case float64:
			hadithNum = int(v)
		case int:
			hadithNum = v
		case string:
			fmt.Sscanf(v, "%d", &hadithNum)
		default:
			continue
		}

		globalID := fmt.Sprintf("%s-%d", bookSlug, hadithNum)

		existingHadith, _ := s.hadithRepo.GetByGlobalID(globalID)
		if existingHadith != nil {
			continue
		}

		hadith := &models.Hadith{
			GlobalID:  globalID,
			BookID:    book.ID,
			ChapterID: 1,
			Number:    hadithNum,
		}

		if err := s.hadithRepo.Create(hadith); err != nil {
			fmt.Printf("Warning: failed to create hadith %s: %v\n", globalID, err)
			continue
		}

		text := &models.HadithText{
			HadithID: hadith.ID,
			Lang:     lang,
			Text:     hadithData.Text,
		}

		if err := s.hadithRepo.CreateText(text); err != nil {
			fmt.Printf("Warning: failed to create text for %s: %v\n", globalID, err)
		}
	}

	return nil
}

func (s *Seeder) SeedChapters(bookSlug, editionName string) error {
	url := fmt.Sprintf("https://cdn.jsdelivr.net/gh/fawazahmed0/hadith-api@1/editions/%s.min.json", editionName)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch chapters: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var edition EditionResponse
	if err := json.Unmarshal(data, &edition); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	book, err := s.bookRepo.GetBySlug(bookSlug)
	if err != nil || book == nil {
		return fmt.Errorf("book not found: %s", bookSlug)
	}

	for chapterNumStr, chapterData := range edition.Metadata.SectionDetail {
		if chapterMap, ok := chapterData.(map[string]interface{}); ok {
			title := ""
			if titleVal, exists := chapterMap["title"]; exists {
				if titleStr, ok := titleVal.(string); ok {
					title = titleStr
				}
			}

			chapterNum, err := strconv.Atoi(chapterNumStr)
			if err != nil {
				fmt.Printf("Warning: invalid chapter number %q: %v\n", chapterNumStr, err)
				continue
			}

			chapter := &models.Chapter{
				BookID:  book.ID,
				Number:  chapterNum,
				TitleAr: title,
				TitleEn: title,
			}

			if err := s.chapterRepo.Create(chapter); err != nil {
				fmt.Printf("Warning: failed to create chapter %d: %v\n", chapterNum, err)
			}
		}
	}

	return nil
}

func (s *Seeder) SeedBooks() error {
	books := []models.Book{
		{Slug: "bukhari", NameAr: "صحيح البخاري", NameEn: "Sahih al-Bukhari", Totals: 7563},
		{Slug: "muslim", NameAr: "صحيح مسلم", NameEn: "Sahih Muslim", Totals: 5037},
		{Slug: "abudawud", NameAr: "سنن أبي داود", NameEn: "Sunan Abu Dawud", Totals: 5274},
		{Slug: "tirmidhi", NameAr: "جامع الترمذي", NameEn: "Jami At-Tirmidhi", Totals: 3956},
		{Slug: "nasai", NameAr: "سنن النسائي", NameEn: "Sunan an-Nasai", Totals: 5761},
		{Slug: "ibnmajah", NameAr: "سنن ابن ماجه", NameEn: "Sunan Ibn Majah", Totals: 4332},
	}

	for _, book := range books {
		existing, _ := s.bookRepo.GetBySlug(book.Slug)
		if existing == nil {
			if err := s.bookRepo.Create(&book); err != nil {
				return fmt.Errorf("failed to create book %s: %w", book.Slug, err)
			}
			fmt.Printf("Created book: %s\n", book.Slug)
		}
	}

	return nil
}