package models

import (
	"github.com/google/uuid"
	"time"
)

type Nomination struct {
	GameID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"game_id"`
	CategoryID     uint      `gorm:"primaryKey" json:"category_id"`
	NominationDate time.Time `gorm:"type:date;not null" json:"nomination_date"`

	Game     Game     `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE" json:"game"`
	Category Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"category"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Nomination) TableName() string {
	return "nominations"
}
