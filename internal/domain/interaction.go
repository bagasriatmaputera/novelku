package domain

import "time"

type LikeTarget string

const (
	LikeNovel   LikeTarget = "novel"
	LikeChapter LikeTarget = "chapter"
	LikeComment LikeTarget = "comment"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	User      User      `json:"user,omitempty"`
	ChapterID uint      `gorm:"index;not null" json:"chapter_id"`
	Chapter   Chapter   `json:"-"`
	ParentID  *uint     `gorm:"index" json:"parent_id"`
	Replies   []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
	Body      string    `gorm:"type:text;not null" json:"body"`
	Timestamp
}

type Like struct {
	UserID    uint       `gorm:"primaryKey"`
	Target    LikeTarget `gorm:"primaryKey;size:20;index:idx_like_target,priority:1"`
	TargetID  uint       `gorm:"primaryKey;index:idx_like_target,priority:2"`
	CreatedAt time.Time  `json:"created_at"`
}

type Library struct {
	UserID    uint  `gorm:"primaryKey"`
	User      User  `json:"-"`
	NovelID   uint  `gorm:"primaryKey;index"`
	Novel     Novel `json:"novel,omitempty"`
	CreatedAt time.Time
}
