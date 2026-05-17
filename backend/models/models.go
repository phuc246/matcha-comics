package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"unique;not null" json:"username"`
	Password  string         `gorm:"not null" json:"-"`
	Role      string         `gorm:"default:editor" json:"role"` // admin, editor
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Story struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Slug        string         `gorm:"unique;not null" json:"slug"`
	Type        string         `gorm:"not null" json:"type"` // comic, novel
	Status      string         `gorm:"default:ongoing" json:"status"` // ongoing, completed, hiatus
	Description string         `json:"description"`
	Author      string         `json:"author"`
	Publisher   string         `json:"publisher"`
	CoverURL    string         `json:"coverUrl"`
	Views       int            `gorm:"default:0" json:"views"`
	Rating      float32        `gorm:"default:0" json:"rating"`
	IsHot       bool           `gorm:"default:false" json:"isHot"`
	Genres      []Genre        `gorm:"many2many:story_genres;" json:"genres"`
	Chapters    []Chapter      `json:"chapters"`
	LatestChapter float64      `gorm:"default:0" json:"latestChapter"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Chapter struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	StoryID   uint            `gorm:"index" json:"storyId"`
	Number    float64         `gorm:"not null" json:"number"`
	Title     string          `json:"title"`
	Note      string          `gorm:"type:text" json:"note"`
	Servers   []ChapterServer `json:"servers"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

type Media struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FileName  string    `json:"fileName"`
	FilePath  string    `json:"filePath"` // Key trên R2 hoặc đường dẫn local
	URL       string    `json:"url"`      // URL truy cập công khai
	Provider  string    `json:"provider"` // "r2" hoặc "local"
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
}

type ChapterServer struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ChapterID uint           `gorm:"index" json:"chapterId"`
	Name      string         `json:"name"` // VD: "Server 1", "Drive"
	Type      string         `json:"type"` // "image" hoặc "text"
	Content   string         `gorm:"type:text" json:"content"` // Danh sách ảnh (JSON) hoặc Nội dung chữ
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Genre struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `gorm:"unique;not null" json:"name"`
	Slug       string `gorm:"unique;not null" json:"slug"`
	BadgeColor string `json:"badgeColor"` // Mã màu cho badge, VD: #ff0000
}
