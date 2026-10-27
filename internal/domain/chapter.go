package domain

import (
	"time"

	"gorm.io/datatypes"
)

type ChapterStatus string

const (
	ChapterDraft     ChapterStatus = "draft"
	ChapterPublished ChapterStatus = "published"
)

type Chapter struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	NovelID     uint           `gorm:"not null;index:idx_chapter_novel_order;priority:1" json:"novel_id"`
	Novel       Novel          `gorm:"constrained;OnDeleted:CASCADE" json:"-"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Content     datatypes.JSON `gorm:"type:jsonb" json:"content"`
	OrderNo     int            `gorm:"not null;index:idx_chapter_novel_order;priority:2" json:"order_no"`
	Status      ChapterStatus  `gorm:"size:20;not null;index;default:draft" json:"status"`
	PublishedAt *time.Time     `json:"published_at"`
	Timestamp
}
