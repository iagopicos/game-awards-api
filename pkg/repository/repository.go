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

// genericRepository implements Repository interface
type genericRepository[T any] struct {
	db *gorm.DB
}

// NewRepository creates a new generic repository
func NewRepository[T any](db *gorm.DB) Repository[T] {
	return &genericRepository[T]{db: db}
}

// Create inserts a new entity
func (r *genericRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

// FindByID retrieves an entity by its ID
func (r *genericRepository[T]) FindByID(id any) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindAll retrieves all entities
func (r *genericRepository[T]) FindAll() ([]T, error) {
	var entities []T
	err := r.db.Find(&entities).Error
	return entities, err
}

// Update updates an existing entity
func (r *genericRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

// Delete removes an entity by ID
func (r *genericRepository[T]) Delete(id any) error {
	var entity T
	return r.db.Delete(&entity, id).Error
}

// CreateTx creates an entity within a transaction
func (r *genericRepository[T]) CreateTx(tx *gorm.DB, entity *T) error {
	return tx.Create(entity).Error
}

// UpdateTx updates an entity within a transaction
func (r *genericRepository[T]) UpdateTx(tx *gorm.DB, entity *T) error {
	return tx.Save(entity).Error
}

// DeleteTx deletes an entity within a transaction
func (r *genericRepository[T]) DeleteTx(tx *gorm.DB, id any) error {
	var entity T
	return tx.Delete(&entity, id).Error
}
