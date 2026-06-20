package models

// BookListResponse wraps a list of books.
type BookListResponse struct {
	Data []Book `json:"data"`
}

// ChapterListResponse wraps a list of chapters.
type ChapterListResponse struct {
	Data []Chapter `json:"data"`
}

// Pagination holds pagination metadata.
type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

// HadithListResponse wraps a paginated list of hadiths.
type HadithListResponse struct {
	Data       []Hadith    `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// ErrorResponse is the standard error response body.
type ErrorResponse struct {
	Error string `json:"error"`
}

// HealthResponse is returned by the health endpoint.
type HealthResponse struct {
	Status string `json:"status"`
}

// SearchResult is one hit: a hadith plus the matching text snippet and language.
type SearchResult struct {
	Hadith *Hadith `json:"hadith"`
	Text   string  `json:"text"`
	Lang   string  `json:"lang"`
}

// SearchListResponse wraps a paginated list of search hits.
type SearchListResponse struct {
	Data       []SearchResult `json:"data"`
	Pagination Pagination     `json:"pagination"`
}