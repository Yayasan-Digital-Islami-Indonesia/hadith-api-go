# Design Spec: Hadith REST API (Go)

- **Date:** 2026-06-12
- **Topic:** High-performance, lightweight Hadith API using Go and SQLite.

## Overview
A production-ready REST API serving the Kutub al-Sittah (6 canonical books) with Arabic text, English, and Indonesian translations. Uses an embedded SQLite database for deployment simplicity and FTS5 for fast search.

## Architecture (Three-Tier)
1.  **API Layer:** Gin-based handlers, JSON serialization, middleware (CORS, Rate Limiting, Logging).
2.  **Service Layer:** Business-to-data mapping, text processing, random selection logic, search orchestration.
3.  **Repository Layer:** GORM with SQLite driver, FTS5 virtual table management, database migrations/seeding.

## Database Schema
- `books`: id, slug, name_ar, name_en, totals.
- `chapters`: id, book_id, number, title_ar, title_en, title_id.
- `hadiths`: id, global_id, book_id, chapter_id, number.
- `hadith_texts`: id, hadith_id, lang, text, narration_chain.
- `hadith_fts` (Virtual): rowid, hadith_id, text_ar, text_en, text_id, chapter_title, book_slug.

## Ingestion Pipeline
1.  Target: `fawazahmed0/hadith-api` GitHub repository.
2.  Mechanism: Go-based CLI tool (seeder) that fetches JSON, parses it, and populates SQLite.
3.  Language: Arabic-first (mandatory), translations added as available.

## API Specification (v1)
- `GET /health`
- `GET /docs` (Swagger UI)
- `GET /api/v1/books`
- `GET /api/v1/books/:id` (id or slug)
- `GET /api/v1/books/:id/chapters`
- `GET /api/v1/books/:id/chapters/:chapter_id`
- `GET /api/v1/hadith/:id` (global id)
- `GET /api/v1/books/:id/hadith/:number`
- `GET /api/v1/search?q=...&page=1&limit=20`
- `GET /api/v1/random`

## Environment Configuration
- `DATABASE_PATH`: Default `./hadith.db`
- `PORT`: Default `8080`
- `LOG_LEVEL`: Default `info`
- `ALLOWED_ORIGINS`: Default `*`
- `DEFAULT_LANGUAGE`: Default `ar`

## Non-Functional Requirements
- **Performance:** Sub-10ms search on Kutub al-Sittah records.
- **Memory:** Optimized GORM queries to maintain low profile.
- **Verification:** Unit tests for services, integration tests for API handlers.
- **Deployment:** Docker & Docker Compose setup.
