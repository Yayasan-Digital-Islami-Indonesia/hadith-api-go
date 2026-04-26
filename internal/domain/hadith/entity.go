package hadith

// HadithText holds the localized text of a hadith.
type HadithText struct {
	Lang string `json:"lang"`
	Text string `json:"text"`
}

// Grade represents an authenticity grade for a hadith.
type Grade struct {
	Name  string `json:"name"`
	Grade string `json:"grade"`
}

// Hadith represents an individual hadith entry.
type Hadith struct {
	Collection   string       `json:"collection"`
	BookNumber   string       `json:"bookNumber"`
	ChapterId    string       `json:"chapterId"`
	HadithNumber string       `json:"hadithNumber"`
	Body         []HadithText `json:"body"`
	Grades       []Grade      `json:"grades"`
}
