package model

import (
	"time"

	"github.com/guregu/null/v5"
	"gorm.io/gorm"
	t "goyave.dev/goyave/v5/util/typeutil"
)

type Article struct {
	Author *User

	CreatedAt t.Undefined[time.Time]      `json:",omitzero"`
	UpdatedAt t.Undefined[null.Time]      `json:",omitzero"`
	DeletedAt t.Undefined[gorm.DeletedAt] `json:",omitzero"`
	Title     t.Undefined[string]         `json:",omitzero"`
	Contents  t.Undefined[string]         `json:",omitzero"`
	Slug      t.Undefined[string]         `json:",omitzero"`
	AuthorID  t.Undefined[int64]          `json:",omitzero"`
	ID        t.Undefined[int64]          `gorm:"primarykey" json:",omitzero"`
}

func (Article) TableName() string {
	return "articles"
}
