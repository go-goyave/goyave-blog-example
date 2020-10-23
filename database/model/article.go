package model

import (
	"time"

	"github.com/System-Glitch/goyave/v3/database"
)

func init() {
	database.RegisterModel(&Article{})
}

// Article represents an article posted by a user.
type Article struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string `gorm:"type:char(200);not null"`
	Contents  string `gorm:"type:longtext;not null"`
	Slug      string `gorm:"type:char(80);not null;unique_index"`
	AuthorID  uint   `json:"-"`
	Author    *User  `gorm:"constraint:OnDelete:CASCADE;" json:",omitempty"`
}
