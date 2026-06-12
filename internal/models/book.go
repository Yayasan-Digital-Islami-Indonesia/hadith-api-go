package models

type Book struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Slug     string    `gorm:"uniqueIndex" json:"slug"`
	NameAr   string    `json:"name_ar"`
	NameEn   string    `json:"name_en"`
	Totals   int       `json:"totals"`
	Chapters []Chapter `gorm:"foreignKey:BookID" json:"-"`
}

func (Book) TableName() string {
	return "books"
}