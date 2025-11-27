package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UserID      uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"user_id"`
	Username    string       `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Name        string       `gorm:"type:varchar(255);not null" json:"name"`
	Email       string       `gorm:"type:varchar(255);not null" json:"email"`
	Password    string       `gorm:"type:varchar(255);not null" json:"password"`
	Predictions []Prediction `gorm:"foreignKey:UserID" json:"predictions,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UdatedAt    time.Time    `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
