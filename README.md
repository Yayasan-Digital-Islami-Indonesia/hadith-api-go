# Hadith API Go

### REST API Hadith Kutub al-Sittah (6 Kitab Induk)

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![SQLite](https://img.shields.io/badge/SQLite-FTS5-07405E?style=flat&logo=sqlite&logoColor=white)](https://www.sqlite.org/)
[![Gin](https://img.shields.io/badge/Gin-Web_Framework-000000?style=flat&logo=go)](https://gin-gonic.com/)
[![License](https://img.shields.io/github/license/ydgi/hadith-api-go?style=flat&colorA=080f12&colorB=1fa669)](LICENSE)

---

REST API untuk data hadith dari 6 kitab induk (Kutub al-Sittah): Sahih Bukhari, Sahih Muslim, Sunan Abu Dawud, Jami At-Tirmidhi, Sunan an-Nasai, dan Sunan Ibn Majah. Menyediakan teks Arab, terjemahan Inggris, dan Indonesia.

- **Cepat** — P95 < 200ms
- **Ringan** — Single binary, SQLite embedded
- **Simple** — JSON response

---

## Quick Start

```bash
git clone https://github.com/ydgi/hadith-api-go.git
cd hadith-api-go
go mod download
make migrate && make seed && make run
```

Server jalan di `http://localhost:8080`

**Docker:**

```bash
docker build -t hadith-api-go -f deploy/Dockerfile .
docker run -p 8080:8080 -e ALLOWED_ORIGINS=https://yourapp.com hadith-api-go
```

---

## Endpoint

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/health` | Health check |
| GET | `/docs` | Swagger UI |
| GET | `/api/v1/books` | Daftar 6 kitab |
| GET | `/api/v1/books/:id` | Detail kitab (id atau slug) |
| GET | `/api/v1/books/:id/chapters` | Daftar bab dalam kitab |
| GET | `/api/v1/books/:id/chapters/:chapter_id` | Hadith dalam bab |
| GET | `/api/v1/hadith/:id` | Hadith by ID |
| GET | `/api/v1/books/:id/hadith/:number` | Hadith by nomor dalam kitab |
| GET | `/api/v1/search?q=...` | Cari hadith by keyword |
| GET | `/api/v1/random` | Hadith acak |

---

## Contoh

**Daftar Kitab:**

```bash
curl http://localhost:8080/api/v1/books
```

**Baca Hadith:**

```bash
curl http://localhost:8080/api/v1/hadith/1
```

**Cari:**

```bash
curl "http://localhost:8080/api/v1/search?q=niat&page=1&limit=10"
```

**Hadith Acak:**

```bash
curl http://localhost:8080/api/v1/random
```

---

## Query Parameters

| Param | Value |
|-------|-------|
| `q` | Keyword pencarian (required untuk `/search`) |
| `page` | Halaman (default: `1`) |
| `limit` | Jumlah per halaman (default: `20`, max: `100`) |

---

## Konfigurasi

| Env Variable | Default | Deskripsi |
|--------------|---------|-----------|
| `DATABASE_PATH` | `./hadith.db` | Path database SQLite |
| `PORT` | `8080` | Port server |
| `ALLOWED_ORIGINS` | `*` | CORS allowed origins |
| `LOG_LEVEL` | `info` | Level logging (debug/info/warn/error) |

---

## Tech Stack

```
Go 1.24+ • Gin • GORM • SQLite FTS5 • Swagger/OpenAPI
```

---

## Development

```bash
make build   # build binaries
make run     # jalankan server
make test    # run tests
make seed    # seed data dari hadith-api
make clean   # hapus bin/ dan database
```

---

## Kontribusi via Fork

```bash
# 1. Fork repo, lalu clone fork kamu
git clone https://github.com/YOUR_USERNAME/hadith-api-go.git
cd hadith-api-go

# 2. Tambah upstream
git remote add upstream https://github.com/ydgi/hadith-api-go.git

# 3. Buat branch, coding, test
git checkout -b feature/fitur-kamu
# ... edit code ...
make test && make build

# 4. Push ke fork, buat PR
git push origin feature/fitur-kamu
# Buka PR di GitHub → "Compare across forks"
```

Lihat [CONTRIBUTING.md](CONTRIBUTING.md) untuk detail.

---

## Data Source

Data hadith diambil dari [fawazahmed0/hadith-api](https://github.com/fawazahmed0/hadith-api) via jsDelivr CDN.

---

## License

MIT
