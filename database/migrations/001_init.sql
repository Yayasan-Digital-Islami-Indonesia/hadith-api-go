CREATE TABLE IF NOT EXISTS books (
  id INTEGER PRIMARY KEY,
  slug TEXT UNIQUE NOT NULL,
  name_ar TEXT NOT NULL,
  name_en TEXT NOT NULL,
  totals INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS chapters (
  id INTEGER PRIMARY KEY,
  book_id INTEGER NOT NULL,
  number INTEGER NOT NULL,
  title_ar TEXT NOT NULL,
  title_en TEXT NOT NULL,
  title_id TEXT,
  FOREIGN KEY(book_id) REFERENCES books(id),
  UNIQUE(book_id, number)
);

CREATE TABLE IF NOT EXISTS hadiths (
  id INTEGER PRIMARY KEY,
  global_id TEXT UNIQUE NOT NULL,
  book_id INTEGER NOT NULL,
  chapter_id INTEGER NOT NULL,
  number INTEGER NOT NULL,
  FOREIGN KEY(book_id) REFERENCES books(id),
  FOREIGN KEY(chapter_id) REFERENCES chapters(id),
  UNIQUE(book_id, number)
);

CREATE TABLE IF NOT EXISTS hadith_texts (
  id INTEGER PRIMARY KEY,
  hadith_id INTEGER NOT NULL,
  lang TEXT NOT NULL,
  text TEXT NOT NULL,
  narration_chain TEXT,
  FOREIGN KEY(hadith_id) REFERENCES hadiths(id),
  UNIQUE(hadith_id, lang)
);

CREATE VIRTUAL TABLE IF NOT EXISTS hadith_fts USING fts5(
  hadith_id,
  text_ar,
  text_en,
  text_id,
  chapter_title,
  book_slug
);

CREATE INDEX IF NOT EXISTS idx_chapters_book ON chapters(book_id);
CREATE INDEX IF NOT EXISTS idx_hadiths_book ON hadiths(book_id);
CREATE INDEX IF NOT EXISTS idx_hadiths_chapter ON hadiths(chapter_id);
CREATE INDEX IF NOT EXISTS idx_hadith_texts_hadith ON hadith_texts(hadith_id);