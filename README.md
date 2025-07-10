# Web Scraper Backend

A backend service in Go (Golang) that analyzes a given website URL and returns key information about the page, with caching using MySQL and GORM.

## Features
- Accepts a website URL via a REST API
- Crawls and analyzes the page for:
  - HTML version
  - Page title
  - Count of heading tags by level (H1, H2, etc.)
  - Number of internal vs. external links
  - Number of inaccessible links (4xx or 5xx)
  - Presence of a login form
- Caches analysis results in MySQL for 24 hours (per URL)

---

## Project Structure

```
backend/
├── main.go                # Entry point, loads .env, starts Gin server
├── types/                 # Data types and GORM models
│   └── analyze.go         # API request/response structs
│   └── analysis.go        # GORM model for caching (Analysis struct)
├── routes/                # Gin route setup
│   └── routes.go
├── handlers/              # HTTP handler functions
│   └── analyze.go
├── helpers/               # Helper functions (DB, HTML parsing, caching)
│   └── db.go              # DB init, caching logic, SHA-256 URL hashing
├── .env                   # Environment variables (not committed)
├── go.mod, go.sum         # Go modules
```

---

## Go Packages Used

- [`github.com/gin-gonic/gin`](https://github.com/gin-gonic/gin): Web framework for building REST APIs.
- [`github.com/PuerkitoBio/goquery`](https://github.com/PuerkitoBio/goquery): HTML parsing and scraping (jQuery-like syntax).
- [`gorm.io/gorm`](https://gorm.io/): ORM for Go, used for MySQL database access and migrations.
- [`gorm.io/driver/mysql`](https://gorm.io/docs/connecting_to_the_database.html#MySQL): MySQL driver for GORM.
- [`github.com/joho/godotenv`](https://github.com/joho/godotenv): Loads environment variables from a `.env` file.

---

## Caching Mechanism

- **Database:** MySQL table stores analysis results for each URL.
- **Uniqueness:** URLs are hashed using SHA-256 and stored in a `url_hash` column (binary, 32 bytes, unique index) to avoid index length issues.
- **TTL:** Each cache entry has an `analyzed_at` timestamp. Cached results are valid for 24 hours. After that, a new analysis is performed and the cache is updated.
- **ORM:** All DB operations use GORM, including upserts (insert or update on conflict).

---

## Setup & Running

### 1. Clone the repository
```
git clone https://github.com/Azer-ch/web-scraper-backend 
cd web-scraper-backend
```

### 2. Install Go dependencies
```
go mod tidy
```

### 3. Set up MySQL
- Create a database (e.g., `webscraper`).
- Ensure you have a user and password with access.

### 4. Create a `.env` file in the project root:
```
MYSQL_DSN=user:password@tcp(127.0.0.1:3306)/webscraper?parseTime=true
```
Replace `user` and `password` with your MySQL credentials.

### 5. Run the server
```
go run main.go
```
The server will start on `localhost:8080` by default.

---

## API Usage

### Analyze a URL
**POST** `/analyze`

**Request Body:**
```json
{
  "url": "https://example.com"
}
```

**Response Example:**
```json
{
  "html_version": "HTML5",
  "title": "Example Domain",
  "headings": {"h1": 1, "h2": 0, ...},
  "internal_links": 2,
  "external_links": 3,
  "inaccessible_links": 1,
  "has_login_form": false
}
```

---

## Notes
- The `.env` file is **not** committed to version control (see `.gitignore`).
- GORM auto-migrates the database schema on startup.
- The cache table uses a SHA-256 hash of the URL for uniqueness and efficient lookups.
- If you change the model, you may need to drop and recreate the table in development.

---

## License
MIT 