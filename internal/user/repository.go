package user

import (
	"context"
	"errors"

	"github.com/bagasriatmaputera/novelku/internal/domain"
	"gorm.io/gorm"
)

var (
	errNotFound  = errors.New("User tidak ditemukan")
	errDuplicate = errors.New("User sudah ada")
)

// Interface: service bergantung pada ini, bukan pada GORM, jadi mudah di-mock saat test.
type Repository interface {
	// context from go itu untuk memastikan jika request batal maka query ikut batal
	Create(ctx context.Context, u *domain.User) error
	FindByID(ctx context.Context, id uint) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	ExistingByEmailOrUsername(ctx context.Context, email string, username string) (bool, error)
}

// repository gorm
type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// NB: context from go itu untuk memastikan jika request batal maka query ikut batal

func (r *repository) Create(ctx context.Context, u *domain.User) error {
	err := r.db.WithContext(ctx).Create(u).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return errDuplicate
	}
	return err
}

func (r *repository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errNotFound
	}
	return &u, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, errNotFound) {
		return nil, errNotFound
	}
	return &u, nil
}

// GORM, dengan context ini, pakai tabel milik model User, saring baris yang email atau username-nya cocok,

// hitung jumlahnya dan simpan ke count, lalu ambil error-nya (kalau ada) ke err
func (r *repository) ExistingByEmailOrUsername(ctx context.Context, email string, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.User{}).Where("email = ? OR username = ?", email, username).Count(&count).Error
	return count > 0, err
}
