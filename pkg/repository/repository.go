package repository

import (
	"context"

	"gorm.io/gorm"
)

// Repository defines the generic CRUD operations
type Repository[T any] interface {
	// Regular operations (no transaction)
	Create(ctx context.Context, entity *T) error
	FindByID(ctx context.Context, id any) (*T, error)
	FindAll(ctx context.Context) ([]T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id any) error

	// Transaction variants (for writes only)
	CreateTx(ctx context.Context, tx *gorm.DB, entity *T) error
	UpdateTx(ctx context.Context, tx *gorm.DB, entity *T) error
	DeleteTx(ctx context.Context, tx *gorm.DB, id any) error
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
func (r *genericRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// FindByID retrieves an entity by its ID
func (r *genericRepository[T]) FindByID(ctx context.Context, id any) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindAll retrieves all entities
func (r *genericRepository[T]) FindAll(ctx context.Context) ([]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).Find(&entities).Error
	return entities, err
}

// Update updates an existing entity
func (r *genericRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete removes an entity by ID
func (r *genericRepository[T]) Delete(ctx context.Context, id any) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}

// CreateTx creates an entity within a transaction
func (r *genericRepository[T]) CreateTx(ctx context.Context, tx *gorm.DB, entity *T) error {
	return tx.WithContext(ctx).Create(entity).Error
}

// UpdateTx updates an entity within a transaction
func (r *genericRepository[T]) UpdateTx(ctx context.Context, tx *gorm.DB, entity *T) error {
	return tx.WithContext(ctx).Save(entity).Error
}

// DeleteTx deletes an entity within a transaction
func (r *genericRepository[T]) DeleteTx(ctx context.Context, tx *gorm.DB, id any) error {
	var entity T
	return tx.WithContext(ctx).Delete(&entity, id).Error
}
