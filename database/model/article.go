package model

import (
	"time"

	"github.com/guregu/null/v5"
	"gorm.io/gorm"
)

type Article struct {
	Author *User

	CreatedAt time.Time
	UpdatedAt null.Time
	DeletedAt gorm.DeletedAt
	Title     string
	Contents  string
	Slug      string
	AuthorID  uint
	ID        uint `gorm:"primarykey"`
}

func (Article) TableName() string {
	return "articles"
}
