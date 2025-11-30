package models

import (
	"github.com/google/uuid"
	"time"
)

type Award struct {
	AwardID      uint      `gorm:"primaryKey;autoIncrement" json:"award_id"`
	CategoryID   uint      `gorm:"not null;index" json:"category_id"`
	GameWinnerID uuid.UUID `gorm:"type:uuid;not null;index" json:"game_winner_id"`
	AwardDate    time.Time `gorm:"type:date;not null" json:"award_date"`

	Category   Category  `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT" json:"category"`
	GameWinner Game      `gorm:"foreignKey:GameWinnerID;constraint:OnDelete:RESTRICT" json:"game_winner"`
	CreatedAt  time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}

func (Award) TableName() string {
	return "awards"
}
