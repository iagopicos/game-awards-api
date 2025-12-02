package repository

import (
	"gorm.io/gorm"
)

// Repository defines the generic CRUD operations
type Repository[T any] interface {
	// Regular operations (no transaction)
	Create(entity *T) error
	FindByID(id any) (*T, error)
	FindAll() ([]T, error)
	Update(entity *T) error
	Delete(id any) error

	// Transaction variants (for writes only)
	CreateTx(tx *gorm.DB, entity *T) error
	UpdateTx(tx *gorm.DB, entity *T) error
	DeleteTx(tx *gorm.DB, id any) error
}
