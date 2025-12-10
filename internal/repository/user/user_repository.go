package repository

import (
	"context"

	"github.com/iagopicos/game-awards-api/internal/models"
	"github.com/iagopicos/game-awards-api/pkg/repository"
)

// UserRepository interface

type UserRepository interface {
	//Base Reoisitory with CRUD operations
	repository.Repository[models.User]

	//Domain-spececific queries

	GetByMail(ctx context.Context, email string) (*models.User, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	GetByName(ctx context.Context, name string) (*models.User, error)
}
