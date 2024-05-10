package dto

import (
	"time"

	"github.com/guregu/null/v5"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Article struct {
	Author *User `json:"author,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt null.Time `json:"updatedAt"`
	Title     string    `json:"title"`
	Contents  string    `json:"contents"`
	Slug      string    `json:"slug"`
	AuthorID  uint      `json:"authorID"`
	ID        uint      `json:"id"`
}

type CreateArticle struct {
	Title    string `json:"title"`
	Contents string `json:"contents"`
	AuthorID uint   `json:"authorID"`
}

type UpdateArticle struct {
	Title    typeutil.Undefined[string] `json:"title"`
	Contents typeutil.Undefined[string] `json:"contents"`
}
