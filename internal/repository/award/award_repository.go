package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/iagopicos/game-awards-api/internal/models"
	"github.com/iagopicos/game-awards-api/pkg/repository"
)

// AwardRepository interface

type AwardRepository interface {
	//Base Repository with CRUD operations
	repository.Repository[models.Award]

	//Domain-specific queries

	GetByCategory(ctx context.Context, categoryID uint) (*models.Award, error)
	GetByGame(ctx context.Context, gameID uuid.UUID) ([]models.Award, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.Award, error)
}
