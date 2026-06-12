package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type SearchRepository struct {
	db *gorm.DB
}

func NewSearchRepository(db *gorm.DB) *SearchRepository {
	return &SearchRepository{db: db}
}

type SearchResult struct {
	HadithID  uint
	Text      string
	Lang      string
	ChapterID uint
	BookSlug  string
}

func (r *SearchRepository) Search(query string, limit, offset int) ([]SearchResult, int64, error) {
	var results []SearchResult
	var total int64

	countFTS := `SELECT COUNT(*) FROM hadith_fts WHERE hadith_fts MATCH ?`
	if err := r.db.Raw(countFTS, query).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	if total > 0 {
		ftsSQL := `
		SELECT fts.hadith_id, ht.text, ht.lang, c.id as chapter_id, b.slug as book_slug
		FROM hadith_fts fts
		JOIN hadith_texts ht ON fts.rowid = ht.id
		JOIN hadiths h ON fts.hadith_id = h.id
		JOIN chapters c ON h.chapter_id = c.id
		JOIN books b ON h.book_id = b.id
		WHERE hadith_fts MATCH ?
		ORDER BY fts.rowid
		LIMIT ? OFFSET ?
		`

		if err := r.db.Raw(ftsSQL, query, limit, offset).Scan(&results).Error; err != nil {
			return nil, 0, fmt.Errorf("search failed: %w", err)
		}
		return results, total, nil
	}

	likeSQL := `
	SELECT h.id as hadith_id, ht.text, ht.lang, c.id as chapter_id, b.slug as book_slug
	FROM hadiths h
	JOIN hadith_texts ht ON h.id = ht.hadith_id
	JOIN chapters c ON h.chapter_id = c.id
	JOIN books b ON h.book_id = b.id
	WHERE ht.text LIKE ?
	ORDER BY h.id
	LIMIT ? OFFSET ?
	`
	searchQuery := "%" + query + "%"
	if err := r.db.Raw(likeSQL, searchQuery, limit, offset).Scan(&results).Error; err != nil {
		return nil, 0, fmt.Errorf("search failed: %w", err)
	}

	countSQL := `
	SELECT COUNT(DISTINCT h.id)
	FROM hadiths h
	JOIN hadith_texts ht ON h.id = ht.hadith_id
	WHERE ht.text LIKE ?
	`
	if err := r.db.Raw(countSQL, searchQuery).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	return results, total, nil
}