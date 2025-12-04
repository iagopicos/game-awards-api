package models

import "time"

type Category struct {
	CategoryID  uint         `gorm:"primaryKey;autoincrement" json:"category_id"`
	Name        string       `gorm:"type:varchar(255);not null;uniqueIndex" json:"name"`
	Nominations []Nomination `gorm:"foreignKey:CategoryID" json:"nominations,omitempty"`
	Awards      []Award      `gorm:"foreignKey:CategoryID" json:"awards,omitempty"`
	Predictions []Prediction `gorm:"foreignKey:CategoryID" json:"predictions,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (Category) TableName() string {
	return "categories"
}
