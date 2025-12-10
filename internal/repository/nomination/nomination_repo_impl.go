package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/iagopicos/game-awards-api/internal/models"
	"github.com/iagopicos/game-awards-api/pkg/repository"
	"gorm.io/gorm"
)

type nominationRepository struct {
	repository.Repository[models.Nomination]
	db *gorm.DB
}

func NewNominationRepository(db *gorm.DB) NominationRepository {
	return &nominationRepository{
		Repository: repository.NewRepository[models.Nomination](db),
		db:         db,
	}
}

func (r *nominationRepository) GetByGameAndCategory(ctx context.Context, gameID uuid.UUID, categoryID uint) (*models.Nomination, error) {
	var nomination models.Nomination
	err := r.db.WithContext(ctx).Where("game_id = ? AND category_id = ?", gameID, categoryID).First(&nomination).Error
	if err != nil {
		return nil, err
	}
	return &nomination, nil
}

func (r *nominationRepository) GetByGame(ctx context.Context, gameID uuid.UUID) ([]models.Nomination, error) {
	var nominations []models.Nomination
	err := r.db.WithContext(ctx).Where("game_id = ?", gameID).Find(&nominations).Error
	if err != nil {
		return nil, err
	}
	return nominations, nil
}

func (r *nominationRepository) GetByCategory(ctx context.Context, categoryID uint) ([]models.Nomination, error) {
	var nominations []models.Nomination
	err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Find(&nominations).Error
	if err != nil {
		return nil, err
	}
	return nominations, nil
}
