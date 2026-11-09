package auth

import "github.com/bagasriatmaputera/novelku/internal/domain"

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30,username"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UserResponse struct {
	ID          uint        `json:"id"`
	Username    string      `json:"username"`
	Email       string      `json:"email"`
	DisplayName string      `json:"display_name"`
	AvatarUrl   string      `json:"avatar_url"`
	Role        domain.Role `json:"role"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"`
}

func toUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		AvatarUrl:   u.AvatarUrl,
		Role:        u.Role,
	}
}
