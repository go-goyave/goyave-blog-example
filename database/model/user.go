package model

import (
	"time"

	"github.com/guregu/null/v5"
	t "goyave.dev/goyave/v5/util/typeutil"
)

type User struct {
	Email     t.Undefined[string] `json:",omitzero"`
	Username  t.Undefined[string] `json:",omitzero"`
	Avatar    null.String         `copier:"-"`
	Password  string
	CreatedAt t.Undefined[time.Time] `json:"createdAt,omitzero"`
	UpdatedAt t.Undefined[null.Time] `json:"updatedAt,omitzero"`
	Articles  []*Article             `gorm:"foreignKey:AuthorID" json:",omitzero"`
	ID        t.Undefined[int64]     `gorm:"primaryKey" json:",omitzero"`
}

func (User) TableName() string {
	return "users"
}
