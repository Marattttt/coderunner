package db

import (
	"context"
	"fmt"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Marattttt/new_new_portfolio/auth/internal/config"
	"github.com/Marattttt/new_new_portfolio/auth/pkg/models"
)

type UsersRepository struct {
	db     *gorm.DB
	logger slog.Logger
}

// TODO: add logger configuration into gorm
func NewUsersRepository(conf *config.DBConfig, logger slog.Logger) (*UsersRepository, error) {
	gormConf := gorm.Config{}
	db, err := gorm.Open(postgres.Open(conf.PostgresURI), &gormConf)

	if err != nil {
		return nil, fmt.Errorf("connecting: %w", err)
	}

	return &UsersRepository{
		db:     db,
		logger: logger,
	}, nil
}

// Populates ID field on success
func (u *UsersRepository) CreateUser(ctx context.Context, user *models.User) error {
	tx := u.db.
		WithContext(ctx).
		Begin()

	defer tx.Commit()

	res := tx.Create(user)
	if res.Error != nil {
		tx.Rollback()
		return fmt.Errorf("db: %w", res.Error)
	}

	return nil
}

func (u UsersRepository) GetUser(ctx context.Context, ID int) (*models.User, error) {
	user := models.User{ID: ID}

	res := u.db.
		WithContext(ctx).
		First(&user)

	if res.Error != nil {
		return nil, fmt.Errorf("db: %w", res.Error)
	}

	return &user, nil
}
