package models

type Hadith struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	GlobalID  string       `gorm:"uniqueIndex" json:"global_id"`
	BookID    uint         `gorm:"index" json:"book_id"`
	ChapterID uint         `gorm:"index" json:"chapter_id"`
	Number    int          `json:"number"`
	Book      Book         `gorm:"foreignKey:BookID" json:"-"`
	Chapter   Chapter      `gorm:"foreignKey:ChapterID" json:"-"`
	Texts     []HadithText `gorm:"foreignKey:HadithID" json:"texts"`
}

type HadithText struct {
	ID             uint    `gorm:"primaryKey" json:"id"`
	HadithID       uint    `gorm:"index" json:"hadith_id"`
	Lang           string  `json:"lang"`
	Text           string  `gorm:"type:text" json:"text"`
	NarrationChain string  `gorm:"type:text" json:"narration_chain"`
	Hadith         Hadith  `gorm:"foreignKey:HadithID" json:"-"`
}

func (Hadith) TableName() string {
	return "hadiths"
}

func (HadithText) TableName() string {
	return "hadith_texts"
}