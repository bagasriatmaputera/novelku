package auth

import (
	"context"
	"errors"
	"time"

	"github.com/bagasriatmaputera/novelku/internal/domain"
	"gorm.io/gorm"
)

var errTokenNotFound = errors.New("refresh token tidak ditemukan")

type TokenRepository interface {
	Create(ctx context.Context, t *domain.RefreshToken) error
	FindByHash(ctx context.Context, hash string) (*domain.RefreshToken, error)
	RevokeByHash(ctx context.Context, hash string) error
	RevokeByAllUser(ctx context.Context, userID uint) error
}

type tokenRepository struct {
	db *gorm.DB
}

func newTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Create(ctx context.Context, t *domain.RefreshToken) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *tokenRepository) FindByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	var t domain.RefreshToken
	err := r.db.WithContext(ctx).Where("token_hash= ?", hash).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errTokenNotFound
	}
	return &t, nil
}
func (r *tokenRepository) RevokeByHash(ctx context.Context, hash string) error {
	return r.db.WithContext(ctx).
		Model(&domain.RefreshToken{}).
		Where("token_hash= ? AND revoke_at is NULL", hash).
		Update("revoke_at", time.Now()).Error
}
func (r *tokenRepository) RevokeByAllUser(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&domain.RefreshToken{}).
		Where("user_id= ? AND revoke_at is NULL", userID).
		Update("revoke_at", time.Now()).Error
}
