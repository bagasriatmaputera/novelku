package domain

import (
	"time"
)

type Role string

const (
	RoleReader Role = "reader"
	RoleAuthor Role = "author"
	RoleAdmin  Role = "admin"
)

type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Username     string `gorm:"size:30;uniqueIndex;not null" json:"username"`
	Email        string `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	DisplayName  string `gorm:"size:100" json:"display_name"`
	Bio          string `gorm:"size:500" json:"bio"`
	AvatarUrl    string `gorm:"size:500" json:"avatar_url"`
	Role         Role   `gorm:"size:20;not null;default:reader;index" json:"role"`
	IsActive     bool   `gorm:"not null;default:true" json:"is_active"`
	Timestamp
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index; not null"`
	User      User      `gorm:"constraint;CASCADE;OnDeleted"`
	TokenHash string    `gorm:"size:64;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	RevokeAt  *time.Time
	CreatedAt time.Time
}
