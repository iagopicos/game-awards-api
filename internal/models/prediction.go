package models

import (
	"time"

	"github.com/google/uuid"
)

type Prediction struct {
	PredictionID    uint      `gorm:"primaryKey;autoIncrement" json:"prediction_id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	CategoryID      uint      `gorm:"not null;index" json:"category_id"`
	GameIDPredicted uuid.UUID `gorm:"type:uuid;not null;index" json:"game_id_predicted"`

	// Relationships
	User          User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Category      Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"category,omitempty"`
	GamePredicted Game     `gorm:"foreignKey:GameIDPredicted;constraint:OnDelete:CASCADE" json:"game_predicted,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Prediction) TableName() string {
	return "predictions"
}
