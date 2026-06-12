package main

import (
	"log"

	"github.com/ydgi/hadith-api-go/internal/config"
	"github.com/ydgi/hadith-api-go/internal/repository"
	"github.com/ydgi/hadith-api-go/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := repository.InitDB(cfg.DatabasePath, cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	seeder := services.NewSeeder(db)

	if err := seeder.SeedBooks(); err != nil {
		log.Fatalf("Failed to seed books: %v", err)
	}

	log.Println("Seeding chapters...")
	for _, ed := range []struct {
		slug    string
		edition string
	}{
		{"bukhari", "ara-bukhari"},
		{"muslim", "ara-muslim"},
		{"abudawud", "ara-abudawud"},
		{"tirmidhi", "ara-tirmidhi"},
		{"nasai", "ara-nasai"},
		{"ibnmajah", "ara-ibnmajah"},
	} {
		if err := seeder.SeedChapters(ed.slug, ed.edition); err != nil {
			log.Printf("Warning: failed to seed chapters for %s: %v", ed.slug, err)
		}
	}

	log.Println("Seeding all 6 books (Arabic, English, Indonesian)...")
	editions := []struct {
		slug    string
		edition string
		lang    string
	}{
		{"bukhari", "ara-bukhari", "ar"},
		{"bukhari", "eng-bukhari", "en"},
		{"bukhari", "ind-bukhari", "id"},
		{"muslim", "ara-muslim", "ar"},
		{"muslim", "eng-muslim", "en"},
		{"muslim", "ind-muslim", "id"},
		{"abudawud", "ara-abudawud", "ar"},
		{"abudawud", "eng-abudawud", "en"},
		{"abudawud", "ind-abudawud", "id"},
		{"tirmidhi", "ara-tirmidhi", "ar"},
		{"tirmidhi", "eng-tirmidhi", "en"},
		{"tirmidhi", "ind-tirmidhi", "id"},
		{"nasai", "ara-nasai", "ar"},
		{"nasai", "eng-nasai", "en"},
		{"nasai", "ind-nasai", "id"},
		{"ibnmajah", "ara-ibnmajah", "ar"},
		{"ibnmajah", "eng-ibnmajah", "en"},
		{"ibnmajah", "ind-ibnmajah", "id"},
	}

	for _, ed := range editions {
		log.Printf("Seeding %s (%s - %s)...", ed.slug, ed.edition, ed.lang)
		if err := seeder.FetchAndSeed(ed.slug, ed.edition, ed.lang); err != nil {
			log.Printf("Failed to seed %s: %v", ed.edition, err)
		}
	}

	log.Println("Seeding complete!")
}
