package models

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	GameID      uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"game_id"`
	Title       string       `gorm:"type:varchar(255);not null" json:"title"`
	ReleaseDate *time.Time   `gorm:"type:date" json:"release_date"`
	ImageURL    string       `gorm:"type:varchar(500)" json:"image_url"`
	Nominations []Nomination `gorm:"foreignKey:GameID" json:"nominations,omitempty"`
	Awards      []Award      `gorm:"foreignKey:GameWinnerID" json:"awards,omitempty"`
	Predictions []Prediction `gorm:"foreignKey:GameIDPredicted" json:"predictions,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (Game) TableName() string {
	return "games"
}
