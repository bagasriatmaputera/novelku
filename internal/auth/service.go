package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/bagasriatmaputera/novelku/internal/domain"
	"github.com/bagasriatmaputera/novelku/internal/user"
	"github.com/bagasriatmaputera/novelku/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists         = errors.New("email atau username sudah terdaftar")
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrAcountDisable      = errors.New("akun dinonaktifkan")
	ErrInvalidToken       = errors.New("refresh token tidak valid")
)

type Service struct {
	users      user.Repository
	tokens     TokenRepository
	jwt        *token.Manager
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewService(users user.Repository, tokens TokenRepository, jwt *token.Manager, accessTTL, refreshTTL time.Duration) *Service {
	return &Service{
		users:      users,
		tokens:     tokens,
		jwt:        jwt,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// membuat lowercase dan menghilangkan spasi awal
	email := strings.ToLower(strings.TrimSpace(req.Email))
	username := strings.TrimSpace(req.Username)

	// melakukan pengecekan existing email dan username
	taken, err := s.users.ExistingByEmailOrUsername(ctx, email, username)
	if err != nil {
		return nil, err
	}

	if taken {
		return nil, ErrUserExists
	}

	// mengacak (hashing) password asli
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create DataBase
	u := &domain.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		DisplayName:  username,
		Role:         domain.RoleReader,
		IsActive:     true,
	}

	if err := s.users.Create(ctx, u); err != nil {
		if errors.Is(err, user.ErrDuplicate) {
			return nil, ErrUserExists
		}
		return nil, err
	}

	return s.issueTokens(ctx, u)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	// Email dibersihkan dari spasi dan diubah lowercase
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Menggunakan repository untuk check user di db
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Compate PasswordHash DB dengan Password Input User
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		return nil, ErrInvalidCredentials
	}

	// Check apakah user aktif
	if !u.IsActive {
		return nil, ErrAcountDisable
	}

	// IssueTokens return AuthResponse dan Error makanya ini tidak error
	return s.issueTokens(ctx, u)
}

func (s *Service) Refresh(ctx context.Context, rawToken string) (*AuthResponse, error) {
	// Melakukan Hash pada token baru
	hash := token.HashRefresh(rawToken)

	// Mecari token hash by hash dan logika err
	rt, err := s.tokens.FindByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, errTokenNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	// Token yang sudah dicabut dipakai lagi: kemungkinan dicuri. Cabut semua sesi user ini.
	if rt.RevokedAt != nil {
		_ = s.tokens.RevokeByAllUser(ctx, rt.ID)
		return nil, ErrInvalidToken
	}
	if time.Now().After(rt.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	// Check user dengan hasil refresh
	u, err := s.users.FindByID(ctx, rt.UserID)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}
	if !u.IsActive {
		return nil, ErrAcountDisable
	}

	// Refresh Token Rotation: Cabut token lama sebelum menerbitkan token baru.
	// Jika token lama dipakai lagi di kemudian hari, sistem tahu ada indikasi token dicuri.
	if err := s.tokens.RevokeByHash(ctx, hash); err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, u)

}

func (s *Service) issueTokens(ctx context.Context, u *domain.User) (*AuthResponse, error) {

	access, err := s.jwt.GenerateAccess(u.ID, string(u.Role))
	if err != nil {
		return nil, err
	}

	// Membuat token baru
	raw, hash, err := token.NewRefresh()
	if err != nil {
		return nil, err
	}

	// Mengisi Refresh Token
	rt := &domain.RefreshToken{
		UserID:    u.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}
	if err := s.tokens.Create(ctx, rt); err != nil {
		return nil, err
	}

	// Return IssueTokens
	return &AuthResponse{
		User:         toUserResponse(u),
		AccessToken:  access,
		RefreshToken: raw,
		ExpiresIn:    int(s.accessTTL.Seconds()),
	}, nil
}
