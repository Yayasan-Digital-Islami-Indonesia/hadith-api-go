package models

type Chapter struct {
	ID      uint      `gorm:"primaryKey" json:"id"`
	BookID  uint      `gorm:"index" json:"book_id"`
	Number  int       `json:"number"`
	TitleAr string    `json:"title_ar"`
	TitleEn string    `json:"title_en"`
	TitleId string    `json:"title_id"`
	Book    Book      `gorm:"foreignKey:BookID" json:"-"`
	Hadiths []Hadith  `gorm:"foreignKey:ChapterID" json:"-"`
}

func (Chapter) TableName() string {
	return "chapters"
}