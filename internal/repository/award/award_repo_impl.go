package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/iagopicos/game-awards-api/internal/models"
	"github.com/iagopicos/game-awards-api/pkg/repository"
	"gorm.io/gorm"
)

type awardRepository struct {
	repository.Repository[models.Award]
	db *gorm.DB
}

func NewAwardRepository(db *gorm.DB) AwardRepository {
	return &awardRepository{
		Repository: repository.NewRepository[models.Award](db),
		db:         db,
	}
}

func (r *awardRepository) GetByCategory(ctx context.Context, categoryID uint) (*models.Award, error) {
	var award models.Award
	err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Order("award_date DESC").First(&award).Error
	if err != nil {
		return nil, err
	}
	return &award, nil
}

func (r *awardRepository) GetByGame(ctx context.Context, gameID uuid.UUID) ([]models.Award, error) {
	var awards []models.Award
	err := r.db.WithContext(ctx).Where("game_winner_id = ?", gameID).Find(&awards).Error
	if err != nil {
		return nil, err
	}
	return awards, nil
}

func (r *awardRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.Award, error) {
	var awards []models.Award
	err := r.db.WithContext(ctx).Where("award_date BETWEEN ? AND ?", startDate, endDate).Find(&awards).Error
	if err != nil {
		return nil, err
	}
	return awards, nil
}
