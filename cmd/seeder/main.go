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

	log.Println("Seeding all 6 books (Arabic)...")
	editions := []struct {
		slug    string
		edition string
	}{
		{"bukhari", "ara-bukhari"},
		{"muslim", "ara-muslim"},
		{"abudawud", "ara-abudawud"},
		{"tirmidhi", "ara-tirmidhi"},
		{"nasai", "ara-nasai"},
		{"ibnmajah", "ara-ibnmajah"},
	}

	for _, ed := range editions {
		log.Printf("Seeding %s (%s)...", ed.slug, ed.edition)
		if err := seeder.FetchAndSeed(ed.slug, ed.edition, "ar"); err != nil {
			log.Printf("Failed to seed %s: %v", ed.slug, err)
		}
	}

	log.Println("Seeding complete!")
}
