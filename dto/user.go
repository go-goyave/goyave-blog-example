package dto

import (
	"time"

	"github.com/guregu/null/v5"
	"goyave.dev/goyave/v5/util/fsutil"
	t "goyave.dev/goyave/v5/util/typeutil"
)

// User the public user DTO. Used to show profiles for example.
type User struct {
	CreatedAt t.Undefined[time.Time] `json:"createdAt,omitzero"`
	UpdatedAt t.Undefined[null.Time] `json:"updatedAt,omitzero"`
	Username  t.Undefined[string]    `json:"username,omitzero"`
	Email     t.Undefined[string]    `json:"email,omitzero"`
	ID        t.Undefined[int64]     `json:"id,omitzero"`
}

// InternalUser contains private user info that should not be exposed to clients.
type InternalUser struct {
	Avatar   null.String `json:"avatar"`
	Password string      `json:"password"`
	User
}

type RegisterUser struct {
	Email    string                                 `json:"email"`
	Username string                                 `json:"username"`
	Password string                                 `json:"password" copier:"-"`
	Avatar   t.Undefined[null.Value[[]fsutil.File]] `json:"avatar" copier:"-"`
}

type UpdateUser struct {
	Email    t.Undefined[string]                    `json:"email"`
	Username t.Undefined[string]                    `json:"username"`
	Password t.Undefined[string]                    `json:"password" copier:"-"`
	Avatar   t.Undefined[null.Value[[]fsutil.File]] `json:"avatar" copier:"-"`
}
