package repository

import (
	"context"
	"database/sql"
	"fmt"

	"hadith-api-go/internal/domain"
	"hadith-api-go/internal/domain/hadith"
)

type HadithRepository struct {
	db *sql.DB
}

func NewHadithRepository(db *sql.DB) hadith.HadithRepository {
	return &HadithRepository{db: db}
}

func (r *HadithRepository) ListByBook(ctx context.Context, collectionName, bookNumber string, limit, offset int) ([]hadith.Hadith, int, error) {
	countRow := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM hadiths WHERE collection_name = ? AND book_number = ?
	`, collectionName, bookNumber)
	var total int
	if err := countRow.Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count hadiths: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT collection_name, book_number, chapter_id, hadith_number
		FROM hadiths
		WHERE collection_name = ? AND book_number = ?
		ORDER BY CAST(hadith_number AS INTEGER)
		LIMIT ? OFFSET ?
	`, collectionName, bookNumber, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list hadiths: %w", err)
	}
	defer rows.Close()

	var hadiths []hadith.Hadith
	for rows.Next() {
		var h hadith.Hadith
		if err := rows.Scan(&h.Collection, &h.BookNumber, &h.ChapterId, &h.HadithNumber); err != nil {
			return nil, 0, fmt.Errorf("scan hadith: %w", err)
		}
		h.Body = r.getHadithBody(ctx, h.Collection, h.HadithNumber)
		h.Grades = r.getHadithGrades(ctx, h.Collection, h.HadithNumber)
		hadiths = append(hadiths, h)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate hadiths: %w", err)
	}

	return hadiths, total, nil
}

func (r *HadithRepository) FindByCollection(ctx context.Context, collectionName, hadithNumber string) (hadith.Hadith, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT collection_name, book_number, chapter_id, hadith_number
		FROM hadiths
		WHERE collection_name = ? AND hadith_number = ?
	`, collectionName, hadithNumber)

	var h hadith.Hadith
	if err := row.Scan(&h.Collection, &h.BookNumber, &h.ChapterId, &h.HadithNumber); err != nil {
		if err == sql.ErrNoRows {
			return hadith.Hadith{}, domain.ErrNotFound
		}
		return hadith.Hadith{}, fmt.Errorf("find hadith: %w", err)
	}
	h.Body = r.getHadithBody(ctx, h.Collection, h.HadithNumber)
	h.Grades = r.getHadithGrades(ctx, h.Collection, h.HadithNumber)
	return h, nil
}

func (r *HadithRepository) Random(ctx context.Context) (hadith.Hadith, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT collection_name, book_number, chapter_id, hadith_number
		FROM hadiths
		ORDER BY RANDOM()
		LIMIT 1
	`)

	var h hadith.Hadith
	if err := row.Scan(&h.Collection, &h.BookNumber, &h.ChapterId, &h.HadithNumber); err != nil {
		if err == sql.ErrNoRows {
			return hadith.Hadith{}, domain.ErrNotFound
		}
		return hadith.Hadith{}, fmt.Errorf("random hadith: %w", err)
	}
	h.Body = r.getHadithBody(ctx, h.Collection, h.HadithNumber)
	h.Grades = r.getHadithGrades(ctx, h.Collection, h.HadithNumber)
	return h, nil
}

func (r *HadithRepository) getHadithBody(ctx context.Context, collectionName, hadithNumber string) []hadith.HadithText {
	rows, err := r.db.QueryContext(ctx, `
		SELECT lang, text FROM hadith_texts WHERE collection_name = ? AND hadith_number = ?
	`, collectionName, hadithNumber)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var texts []hadith.HadithText
	for rows.Next() {
		var t hadith.HadithText
		if err := rows.Scan(&t.Lang, &t.Text); err != nil {
			continue
		}
		texts = append(texts, t)
	}
	return texts
}

func (r *HadithRepository) getHadithGrades(ctx context.Context, collectionName, hadithNumber string) []hadith.Grade {
	rows, err := r.db.QueryContext(ctx, `
		SELECT name, grade FROM hadith_grades WHERE collection_name = ? AND hadith_number = ?
	`, collectionName, hadithNumber)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var grades []hadith.Grade
	for rows.Next() {
		var g hadith.Grade
		if err := rows.Scan(&g.Name, &g.Grade); err != nil {
			continue
		}
		grades = append(grades, g)
	}
	return grades
}
