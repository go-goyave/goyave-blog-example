package dto

import (
	"time"

	"github.com/guregu/null/v5"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Article struct {
	Author *User `json:"author,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitzero"`
	UpdatedAt null.Time `json:"updatedAt,omitzero"`
	Title     string    `json:"title,omitempty"`
	Contents  string    `json:"contents,omitempty"`
	Slug      string    `json:"slug,omitempty"`
	AuthorID  uint      `json:"authorID,omitempty"`
	ID        uint      `json:"id,omitempty"`
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
