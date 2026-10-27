package domain

type NovelStatus string

const (
	NovelDraft     NovelStatus = "draft"
	NovelPublished NovelStatus = "published"
	NovelCompleted NovelStatus = "completed"
)

type Novel struct {
	ID       uint        `gorm:"primaryKey" json:"id"`
	AuthorID uint        `gorm:"size:30;index;not null" json:"author_id"`
	Author   User        `gorm:"foreignKey:AuthorId" json:"author,omitempty"`
	Title    string      `gorm:"size:255;not null" json:"title"`
	Slug     string      `gorm:"size:280;uniqueIndex;not null" json:"slug"`
	Synopsis string      `gorm:"type:text" json:"synopsis"`
	CoverUrl string      `gorm:"size:500" json:"cover_url"`
	Status   NovelStatus `gorm:"size:20;not null;default:draft;index" json:"status"`
	Genres   []Genres    `gorm:"many2many:novel_genres" json:"genres,omitempty"`
	Chapters []Chapter   `json:"chapters,omitempty"`
	Timestamp
}

type Genres struct {
	ID   uint   `gorm:"primaryLey" json:"id"`
	Name string `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Slug string `gorm:"size:80;uniqueIndex;not null" json:"slug"`
}
