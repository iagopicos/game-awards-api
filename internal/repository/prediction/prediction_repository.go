package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/iagopicos/game-awards-api/internal/models"
	"github.com/iagopicos/game-awards-api/pkg/repository"
)

// PredictionRepository interface

type PredictionRepository interface {
	//Base Repository with CRUD operations
	repository.Repository[models.Prediction]

	//Domain-specific queries

	GetByUser(ctx context.Context, userID uuid.UUID) ([]models.Prediction, error)
	GetByUserAndCategory(ctx context.Context, userID uuid.UUID, categoryID uint) (*models.Prediction, error)
	GetByCategory(ctx context.Context, categoryID uint, limit, offset int) ([]models.Prediction, int64, error)
	GetByGame(ctx context.Context, gameID uuid.UUID) ([]models.Prediction, error)
}
