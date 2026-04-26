package repository

import (
	"context"
	"database/sql"
	"fmt"

	"hadith-api-go/internal/domain"
	"hadith-api-go/internal/domain/collection"
)

type CollectionRepository struct {
	db *sql.DB
}

func NewCollectionRepository(db *sql.DB) collection.CollectionRepository {
	return &CollectionRepository{db: db}
}

func (r *CollectionRepository) List(ctx context.Context, limit, offset int) ([]collection.Collection, int, error) {
	countRow := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM collections`)
	var total int
	if err := countRow.Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count collections: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT name, has_books, has_chapters, total_hadith, total_available_hadith
		FROM collections
		ORDER BY id
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list collections: %w", err)
	}
	defer rows.Close()

	var collections []collection.Collection
	for rows.Next() {
		var c collection.Collection
		if err := rows.Scan(&c.Name, &c.HasBooks, &c.HasChapters, &c.TotalHadith, &c.TotalAvailableHadith); err != nil {
			return nil, 0, fmt.Errorf("scan collection: %w", err)
		}
		c.Collection = r.getCollectionNames(ctx, c.Name)
		collections = append(collections, c)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate collections: %w", err)
	}

	return collections, total, nil
}

func (r *CollectionRepository) FindByName(ctx context.Context, name string) (collection.Collection, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT name, has_books, has_chapters, total_hadith, total_available_hadith
		FROM collections
		WHERE name = ?
	`, name)

	var c collection.Collection
	if err := row.Scan(&c.Name, &c.HasBooks, &c.HasChapters, &c.TotalHadith, &c.TotalAvailableHadith); err != nil {
		if err == sql.ErrNoRows {
			return collection.Collection{}, domain.ErrNotFound
		}
		return collection.Collection{}, fmt.Errorf("find collection: %w", err)
	}
	c.Collection = r.getCollectionNames(ctx, c.Name)
	return c, nil
}

func (r *CollectionRepository) getCollectionNames(ctx context.Context, name string) []struct {
	Lang  string `json:"lang"`
	Title string `json:"title"`
	Short string `json:"shortIntro"`
} {
	rows, err := r.db.QueryContext(ctx, `
		SELECT lang, title, short_intro FROM collection_names WHERE collection_name = ?
	`, name)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var names []struct {
		Lang  string `json:"lang"`
		Title string `json:"title"`
		Short string `json:"shortIntro"`
	}
	for rows.Next() {
		var n struct {
			Lang  string `json:"lang"`
			Title string `json:"title"`
			Short string `json:"shortIntro"`
		}
		if err := rows.Scan(&n.Lang, &n.Title, &n.Short); err != nil {
			continue
		}
		names = append(names, n)
	}
	return names
}
