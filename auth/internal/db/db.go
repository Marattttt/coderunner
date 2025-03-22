package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Marattttt/new_new_portfolio/auth/internal/config"
	"github.com/Marattttt/new_new_portfolio/auth/pkg/models"
)

var ErrNotFound = fmt.Errorf("not found")

type DBConn struct{ gorm *gorm.DB }

func Connect(conf *config.DBConfig) (*DBConn, error) {
	gormConf := gorm.Config{}
	db, err := gorm.Open(postgres.Open(conf.PostgresURI), &gormConf)

	if err != nil {
		return nil, fmt.Errorf("connecting: %w", err)
	}

	return &DBConn{db}, nil
}

type UsersRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

// TODO: add logger configuration into gorm
func NewUsersRepository(db *DBConn, logger *slog.Logger) *UsersRepository {
	return &UsersRepository{
		db:     db.gorm,
		logger: logger,
	}
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

func (u *UsersRepository) GetUserEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	res := u.db.Model(&models.User{}).Where("email = ?", email).First(&user)

	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("db: %w", res.Error)
	}

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	return &user, nil
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
