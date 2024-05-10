package dto

import (
	"time"

	"github.com/guregu/null/v5"
	"goyave.dev/goyave/v5/util/fsutil"
	"goyave.dev/goyave/v5/util/typeutil"
)

// User the public user DTO. Used to show profiles for example.
type User struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt null.Time `json:"updatedAt"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	ID        uint      `json:"id"`
}

// InternalUser contains private user info that should not be exposed to clients.
type InternalUser struct {
	Avatar   string `json:"avatar"`
	Password string `json:"password"`
	User
}

type RegisterUser struct {
	Email    string                            `json:"email"`
	Username string                            `json:"username"`
	Password string                            `json:"password" copier:"-"`
	Avatar   typeutil.Undefined[[]fsutil.File] `json:"avatar" copier:"-"`
}

type UpdateUser struct {
	Email    typeutil.Undefined[string]        `json:"email"`
	Username typeutil.Undefined[string]        `json:"username"`
	Password typeutil.Undefined[string]        `json:"password" copier:"-"`
	Avatar   typeutil.Undefined[[]fsutil.File] `json:"avatar" copier:"-"`
}
