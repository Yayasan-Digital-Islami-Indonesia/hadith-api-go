package seed

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
)

// Run loads hadith seed data from dataDir into the database.
// Seed files should be placed in dataDir as JSON or SQL files.
// TODO: implement actual seed loading from data files.
func Run(ctx context.Context, db *sql.DB, dataDir string) error {
	log.Info().Str("dataDir", dataDir).Msg("starting seed")

	if err := seedCollections(ctx, db); err != nil {
		return fmt.Errorf("seed collections: %w", err)
	}

	log.Info().Msg("seed completed")
	return nil
}

// seedCollections inserts placeholder hadith collection data.
// Replace this with actual data loading logic.
func seedCollections(ctx context.Context, db *sql.DB) error {
	collections := []struct {
		name        string
		hasBooks    bool
		hasChapters bool
		totalHadith int
	}{
		{"bukhari", true, true, 7563},
		{"muslim", true, true, 5362},
		{"abudawud", true, true, 5274},
		{"tirmidhi", true, true, 3956},
		{"nasai", true, true, 5761},
		{"ibnmajah", true, true, 4341},
		{"malik", true, true, 1857},
		{"riyadussalihin", true, true, 1905},
	}

	for _, c := range collections {
		_, err := db.ExecContext(ctx, `
			INSERT OR IGNORE INTO collections (name, has_books, has_chapters, total_hadith, total_available_hadith)
			VALUES (?, ?, ?, ?, ?)
		`, c.name, c.hasBooks, c.hasChapters, c.totalHadith, 0)
		if err != nil {
			return fmt.Errorf("insert collection %s: %w", c.name, err)
		}

		// Insert English and Arabic names as placeholders
		names := []struct {
			lang  string
			title string
			short string
		}{
			{"en", collectionTitle(c.name, "en"), ""},
			{"ar", collectionTitle(c.name, "ar"), ""},
		}
		for _, n := range names {
			_, err := db.ExecContext(ctx, `
				INSERT OR IGNORE INTO collection_names (collection_name, lang, title, short_intro)
				VALUES (?, ?, ?, ?)
			`, c.name, n.lang, n.title, n.short)
			if err != nil {
				return fmt.Errorf("insert collection name %s/%s: %w", c.name, n.lang, err)
			}
		}
	}

	return nil
}

func collectionTitle(name, lang string) string {
	titles := map[string]map[string]string{
		"bukhari":        {"en": "Sahih al-Bukhari", "ar": "صحيح البخاري"},
		"muslim":         {"en": "Sahih Muslim", "ar": "صحيح مسلم"},
		"abudawud":       {"en": "Sunan Abu Dawud", "ar": "سنن أبي داود"},
		"tirmidhi":       {"en": "Jami` at-Tirmidhi", "ar": "جامع الترمذي"},
		"nasai":          {"en": "Sunan an-Nasa'i", "ar": "سنن النسائي"},
		"ibnmajah":       {"en": "Sunan Ibn Majah", "ar": "سنن ابن ماجه"},
		"malik":          {"en": "Muwatta Malik", "ar": "موطأ مالك"},
		"riyadussalihin": {"en": "Riyad as-Salihin", "ar": "رياض الصالحين"},
	}
	if t, ok := titles[name]; ok {
		if v, ok := t[lang]; ok {
			return v
		}
	}
	return name
}
