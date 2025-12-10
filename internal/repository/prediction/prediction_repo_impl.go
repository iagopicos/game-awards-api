package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/iagopicos/game-awards-api/internal/models"
	"github.com/iagopicos/game-awards-api/pkg/repository"
	"gorm.io/gorm"
)

type predictionRepository struct {
	repository.Repository[models.Prediction]
	db *gorm.DB
}

func NewPredictionRepository(db *gorm.DB) PredictionRepository {
	return &predictionRepository{
		Repository: repository.NewRepository[models.Prediction](db),
		db:         db,
	}
}

func (r *predictionRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]models.Prediction, error) {
	var predictions []models.Prediction
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&predictions).Error
	if err != nil {
		return nil, err
	}
	return predictions, nil
}

func (r *predictionRepository) GetByUserAndCategory(ctx context.Context, userID uuid.UUID, categoryID uint) (*models.Prediction, error) {
	var prediction models.Prediction
	err := r.db.WithContext(ctx).Where("user_id = ? AND category_id = ?", userID, categoryID).First(&prediction).Error
	if err != nil {
		return nil, err
	}
	return &prediction, nil
}

func (r *predictionRepository) GetByCategory(ctx context.Context, categoryID uint, limit, offset int) ([]models.Prediction, int64, error) {
	var predictions []models.Prediction
	var total int64

	err := r.db.WithContext(ctx).Model(&models.Prediction{}).Where("category_id = ?", categoryID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).Where("category_id = ?", categoryID).Limit(limit).Offset(offset).Find(&predictions).Error
	if err != nil {
		return nil, 0, err
	}

	return predictions, total, nil
}

func (r *predictionRepository) GetByGame(ctx context.Context, gameID uuid.UUID) ([]models.Prediction, error) {
	var predictions []models.Prediction
	err := r.db.WithContext(ctx).Where("game_id_predicted = ?", gameID).Find(&predictions).Error
	if err != nil {
		return nil, err
	}
	return predictions, nil
}
