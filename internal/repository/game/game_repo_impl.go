package repository

import (
	"context"

	"github.com/iagopicos/game-awards-api/internal/models"
	"github.com/iagopicos/game-awards-api/pkg/repository"
	"gorm.io/gorm"
)

type gameRepository struct {
	repository.Repository[models.Game]
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) GameRepository {
	return &gameRepository{
		Repository: repository.NewRepository[models.Game](db),
		db:         db,
	}

}

func (r *gameRepository) GetByDeveloper(ctx context.Context, developer string, limit, offset int) ([]models.Game, int64, error) {
	var games []models.Game
	var total int64
	err := r.db.WithContext(ctx).
		Where("developer = ?", developer).
		Limit(limit).
		Offset(offset).
		Find(&games).Error
	if err := r.db.WithContext(ctx).
		Model(&models.Game{}).
		Where("developer = ?", developer).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return games, total, err
}
func (r *gameRepository) GetByName(ctx context.Context, developer string) (*models.Game, error) {
	return nil, nil
}
