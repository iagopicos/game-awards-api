package repository

import (
	"context"

	"github.com/iagopicos/game-awards-api/internal/models"
	"github.com/iagopicos/game-awards-api/pkg/repository"
)

// GameRepository interface

type GameRepository interface {
	//Base Reoisitory with CRUD operations
	repository.Repository[models.Game]

	//Domain-spececific queries

	GetByDeveloper(ctx context.Context, email string) ([]*models.Game, error)
	GetByName(ctx context.Context, name string) (*models.Game, error)
}
