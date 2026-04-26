package book

// BookName holds the localized name of a book.
type BookName struct {
	Lang string `json:"lang"`
	Name string `json:"name"`
}

// Book represents a book (chapter grouping) within a hadith collection.
type Book struct {
	BookNumber        string     `json:"bookNumber"`
	Book              []BookName `json:"book"`
	HadithStartNumber int        `json:"hadithStartNumber"`
	HadithEndNumber   int        `json:"hadithEndNumber"`
	NumberOfHadiths   int        `json:"numberOfHadiths"`
}
