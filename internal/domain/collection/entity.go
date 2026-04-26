package collection

// Collection represents a hadith collection (e.g., Sahih Bukhari, Sahih Muslim).
type Collection struct {
	Name        string `json:"name"`
	HasBooks    bool   `json:"hasBooks"`
	HasChapters bool   `json:"hasChapters"`
	Collection  []struct {
		Lang  string `json:"lang"`
		Title string `json:"title"`
		Short string `json:"shortIntro"`
	} `json:"collection"`
	TotalHadith          int `json:"totalHadith"`
	TotalAvailableHadith int `json:"totalAvailableHadith"`
}

// CollectionSummary is a compact representation used in list responses.
type CollectionSummary struct {
	Name        string `json:"name"`
	HasBooks    bool   `json:"hasBooks"`
	HasChapters bool   `json:"hasChapters"`
	Collection  []struct {
		Lang  string `json:"lang"`
		Title string `json:"title"`
		Short string `json:"shortIntro"`
	} `json:"collection"`
	TotalHadith          int `json:"totalHadith"`
	TotalAvailableHadith int `json:"totalAvailableHadith"`
}
