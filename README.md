# Hadith API Go

REST API for Hadith collections. Provides 8+ major hadith collections with Arabic and English text.

- **Cepat** — P95 < 200ms
- **Ringan** — Single binary, SQLite embedded
- **Simple** — JSON response

API format inspired by [Sunnah.com API](https://sunnah.stoplight.io/docs/api/).

---

## Quick Start

```bash
git clone https://github.com/alann-maulana/hadith-api-go.git
cd hadith-api-go
go mod download
make migrate && make seed && make run
```

Server berjalan di `http://localhost:8080`

**Docker:**

```bash
docker build -t hadith-api-go .
docker run -p 8080:8080 -e ALLOWED_ORIGINS=https://yourapp.com hadith-api-go
```

---

## Endpoint

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/health` | Health check |
| GET | `/health/ready` | Readiness check |
| GET | `/v1/collections` | Daftar semua koleksi hadist |
| GET | `/v1/collections/:name` | Detail koleksi hadist |
| GET | `/v1/collections/:name/books` | Daftar buku dalam koleksi |
| GET | `/v1/collections/:name/books/:bookNumber` | Detail buku |
| GET | `/v1/collections/:name/books/:bookNumber/hadiths` | Daftar hadist dalam buku |
| GET | `/v1/collections/:name/hadiths/:hadithNumber` | Hadist spesifik dalam koleksi |
| GET | `/v1/hadiths/random` | Hadist acak |
| GET | `/docs` | Dokumentasi API |

---

## Koleksi yang Tersedia

| Name | Judul |
|------|-------|
| `bukhari` | Sahih al-Bukhari |
| `muslim` | Sahih Muslim |
| `abudawud` | Sunan Abu Dawud |
| `tirmidhi` | Jami` at-Tirmidhi |
| `nasai` | Sunan an-Nasa'i |
| `ibnmajah` | Sunan Ibn Majah |
| `malik` | Muwatta Malik |
| `riyadussalihin` | Riyad as-Salihin |

---

## Contoh

**Daftar Koleksi:**
```bash
curl http://localhost:8080/v1/collections
```

**Detail Koleksi:**
```bash
curl http://localhost:8080/v1/collections/bukhari
```

**Daftar Buku:**
```bash
curl http://localhost:8080/v1/collections/bukhari/books
```

**Daftar Hadist dalam Buku:**
```bash
curl "http://localhost:8080/v1/collections/bukhari/books/1/hadiths?page=1&limit=10"
```

**Hadist Spesifik:**
```bash
curl http://localhost:8080/v1/collections/bukhari/hadiths/1
```

**Hadist Acak:**
```bash
curl http://localhost:8080/v1/hadiths/random
```

---

## Query Parameters

| Param | Deskripsi |
|-------|-----------|
| `page` | Nomor halaman (default: `1`) |
| `limit` | Jumlah per halaman (default: `20`, max: `100`) |

---

## Response Format

**Success:**
```json
{
  "data": { ... },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

**Paginated:**
```json
{
  "data": [...],
  "total": 100,
  "page": 1,
  "limit": 20,
  "timestamp": "2024-01-01T00:00:00Z"
}
```

**Error:**
```json
{
  "error": "resource not found",
  "code": "not found",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

---

## Konfigurasi

| Env Variable | Default |
|-------------|---------|
| `DB_PATH` | `./data/hadith.db` |
| `SERVER_PORT` | `8080` |
| `SERVER_HOST` | `0.0.0.0` |
| `ALLOWED_ORIGINS` | _(kosong)_ |
| `APP_VERSION` | `1.0.0` |
| `LOG_LEVEL` | `info` |

---

## Tech Stack

```
Go 1.25+ • Gin • SQLite • Goose • Zerolog
```

---

## Development

```bash
make test    # run tests
make lint    # static analysis
```

## License

MIT