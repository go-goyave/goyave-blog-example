package model

import (
	"time"

	"github.com/guregu/null/v5"
)

type User struct {
	Email     string
	Username  string
	Avatar    null.String `copier:"-"`
	Password  string
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt null.Time  `json:"updatedAt"`
	Articles  []*Article `gorm:"foreignKey:AuthorID"`
	ID        uint       `gorm:"primaryKey"`
}

func (User) TableName() string {
	return "users"
}
