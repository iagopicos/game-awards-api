package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/iagopicos/game-awards-api/internal/models"
	"github.com/iagopicos/game-awards-api/pkg/repository"
)

// NominationRepository interface

type NominationRepository interface {
	//Base Repository with CRUD operations
	repository.Repository[models.Nomination]

	//Domain-specific queries

	GetByGameAndCategory(ctx context.Context, gameID uuid.UUID, categoryID uint) (*models.Nomination, error)
	GetByGame(ctx context.Context, gameID uuid.UUID) ([]models.Nomination, error)
	GetByCategory(ctx context.Context, categoryID uint) ([]models.Nomination, error)
}
