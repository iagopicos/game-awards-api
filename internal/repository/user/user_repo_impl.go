package repository

import (
	"context"

	"github.com/iagopicos/game-awards-api/internal/models"
	"github.com/iagopicos/game-awards-api/pkg/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	repository.Repository[models.User]
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		Repository: repository.NewRepository[models.User](db),
		db:         db,
	}
}

func (r *userRepository) GetByMail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("login = ?", login).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) GetByName(ctx context.Context, name string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
