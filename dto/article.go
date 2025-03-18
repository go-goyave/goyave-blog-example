package dto

import (
	"time"

	"github.com/guregu/null/v5"
	t "goyave.dev/goyave/v5/util/typeutil"
)

type Article struct {
	Author *User `json:"author,omitempty"`

	CreatedAt t.Undefined[time.Time] `json:"createdAt,omitzero"`
	UpdatedAt t.Undefined[null.Time] `json:"updatedAt,omitzero"`
	Title     t.Undefined[string]    `json:"title,omitzero"`
	Contents  t.Undefined[string]    `json:"contents,omitzero"`
	Slug      t.Undefined[string]    `json:"slug,omitzero"`
	AuthorID  t.Undefined[int64]     `json:"authorID,omitzero"`
	ID        t.Undefined[int64]     `json:"id,omitzero"`
}

type CreateArticle struct {
	Title    string `json:"title"`
	Contents string `json:"contents"`
	AuthorID int64  `json:"authorID"`
}

type UpdateArticle struct {
	Title    t.Undefined[string] `json:"title"`
	Contents t.Undefined[string] `json:"contents"`
}
