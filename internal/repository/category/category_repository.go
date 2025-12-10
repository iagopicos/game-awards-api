package repository

import (
	"context"

	"github.com/iagopicos/game-awards-api/internal/models"
	"github.com/iagopicos/game-awards-api/pkg/repository"
)

// CategoryRepository interface

type CategoryRepository interface {
	//Base Repository with CRUD operations
	repository.Repository[models.Category]

	//Domain-specific queries

	GetByName(ctx context.Context, name string) (*models.Category, error)
	GetAllWithNominations(ctx context.Context) ([]models.Category, error)
}
