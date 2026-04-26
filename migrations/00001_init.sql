-- +goose Up
CREATE TABLE IF NOT EXISTS collections (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	has_books INTEGER NOT NULL DEFAULT 1,
	has_chapters INTEGER NOT NULL DEFAULT 1,
	total_hadith INTEGER NOT NULL DEFAULT 0,
	total_available_hadith INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS collection_names (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	collection_name TEXT NOT NULL,
	lang TEXT NOT NULL,
	title TEXT NOT NULL,
	short_intro TEXT NOT NULL DEFAULT '',
	FOREIGN KEY (collection_name) REFERENCES collections(name),
	UNIQUE (collection_name, lang)
);

CREATE TABLE IF NOT EXISTS books (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	collection_name TEXT NOT NULL,
	book_number TEXT NOT NULL,
	hadith_start_number INTEGER NOT NULL DEFAULT 0,
	hadith_end_number INTEGER NOT NULL DEFAULT 0,
	number_of_hadiths INTEGER NOT NULL DEFAULT 0,
	FOREIGN KEY (collection_name) REFERENCES collections(name),
	UNIQUE (collection_name, book_number)
);

CREATE TABLE IF NOT EXISTS book_names (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	collection_name TEXT NOT NULL,
	book_number TEXT NOT NULL,
	lang TEXT NOT NULL,
	name TEXT NOT NULL,
	FOREIGN KEY (collection_name) REFERENCES collections(name),
	UNIQUE (collection_name, book_number, lang)
);

CREATE TABLE IF NOT EXISTS hadiths (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	collection_name TEXT NOT NULL,
	book_number TEXT NOT NULL,
	chapter_id TEXT NOT NULL DEFAULT '',
	hadith_number TEXT NOT NULL,
	FOREIGN KEY (collection_name) REFERENCES collections(name),
	UNIQUE (collection_name, hadith_number)
);

CREATE TABLE IF NOT EXISTS hadith_texts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	collection_name TEXT NOT NULL,
	hadith_number TEXT NOT NULL,
	lang TEXT NOT NULL,
	text TEXT NOT NULL,
	FOREIGN KEY (collection_name) REFERENCES collections(name),
	UNIQUE (collection_name, hadith_number, lang)
);

CREATE TABLE IF NOT EXISTS hadith_grades (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	collection_name TEXT NOT NULL,
	hadith_number TEXT NOT NULL,
	name TEXT NOT NULL,
	grade TEXT NOT NULL,
	FOREIGN KEY (collection_name) REFERENCES collections(name)
);

CREATE INDEX IF NOT EXISTS idx_collection_names_collection ON collection_names (collection_name);
CREATE INDEX IF NOT EXISTS idx_books_collection ON books (collection_name);
CREATE INDEX IF NOT EXISTS idx_book_names_collection ON book_names (collection_name, book_number);
CREATE INDEX IF NOT EXISTS idx_hadiths_collection ON hadiths (collection_name);
CREATE INDEX IF NOT EXISTS idx_hadiths_book ON hadiths (collection_name, book_number);
CREATE INDEX IF NOT EXISTS idx_hadith_texts_hadith ON hadith_texts (collection_name, hadith_number);
CREATE INDEX IF NOT EXISTS idx_hadith_grades_hadith ON hadith_grades (collection_name, hadith_number);

-- +goose Down
DROP INDEX IF EXISTS idx_hadith_grades_hadith;
DROP INDEX IF EXISTS idx_hadith_texts_hadith;
DROP INDEX IF EXISTS idx_hadiths_book;
DROP INDEX IF EXISTS idx_hadiths_collection;
DROP INDEX IF EXISTS idx_book_names_collection;
DROP INDEX IF EXISTS idx_books_collection;
DROP INDEX IF EXISTS idx_collection_names_collection;

DROP TABLE IF EXISTS hadith_grades;
DROP TABLE IF EXISTS hadith_texts;
DROP TABLE IF EXISTS hadiths;
DROP TABLE IF EXISTS book_names;
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS collection_names;
DROP TABLE IF EXISTS collections;
