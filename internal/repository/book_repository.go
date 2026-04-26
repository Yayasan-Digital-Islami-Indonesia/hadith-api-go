package repository

import (
	"context"
	"database/sql"
	"fmt"

	"hadith-api-go/internal/domain"
	"hadith-api-go/internal/domain/book"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) book.BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) List(ctx context.Context, collectionName string, limit, offset int) ([]book.Book, int, error) {
	countRow := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM books WHERE collection_name = ?
	`, collectionName)
	var total int
	if err := countRow.Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count books: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT book_number, hadith_start_number, hadith_end_number, number_of_hadiths
		FROM books
		WHERE collection_name = ?
		ORDER BY CAST(book_number AS INTEGER)
		LIMIT ? OFFSET ?
	`, collectionName, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list books: %w", err)
	}
	defer rows.Close()

	var books []book.Book
	for rows.Next() {
		var b book.Book
		if err := rows.Scan(&b.BookNumber, &b.HadithStartNumber, &b.HadithEndNumber, &b.NumberOfHadiths); err != nil {
			return nil, 0, fmt.Errorf("scan book: %w", err)
		}
		b.Book = r.getBookNames(ctx, collectionName, b.BookNumber)
		books = append(books, b)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate books: %w", err)
	}

	return books, total, nil
}

func (r *BookRepository) FindByNumber(ctx context.Context, collectionName, bookNumber string) (book.Book, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT book_number, hadith_start_number, hadith_end_number, number_of_hadiths
		FROM books
		WHERE collection_name = ? AND book_number = ?
	`, collectionName, bookNumber)

	var b book.Book
	if err := row.Scan(&b.BookNumber, &b.HadithStartNumber, &b.HadithEndNumber, &b.NumberOfHadiths); err != nil {
		if err == sql.ErrNoRows {
			return book.Book{}, domain.ErrNotFound
		}
		return book.Book{}, fmt.Errorf("find book: %w", err)
	}
	b.Book = r.getBookNames(ctx, collectionName, b.BookNumber)
	return b, nil
}

func (r *BookRepository) getBookNames(ctx context.Context, collectionName, bookNumber string) []book.BookName {
	rows, err := r.db.QueryContext(ctx, `
		SELECT lang, name FROM book_names WHERE collection_name = ? AND book_number = ?
	`, collectionName, bookNumber)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var names []book.BookName
	for rows.Next() {
		var n book.BookName
		if err := rows.Scan(&n.Lang, &n.Name); err != nil {
			continue
		}
		names = append(names, n)
	}
	return names
}
