package model

import (
	"reflect"
	"time"

	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/config"
	"github.com/System-Glitch/goyave/v3/database"
	"github.com/System-Glitch/goyave/v3/middleware/ratelimiter"
	"github.com/bxcodec/faker/v3"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

func init() {
	database.RegisterModel(&User{})

	config.Register("app.bcryptCost", config.Entry{
		Value:            10,
		Type:             reflect.Int,
		IsSlice:          false,
		AuthorizedValues: []interface{}{},
	})
}

// User represents a user.
type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string      `gorm:"type:char(100);unique;uniqueIndex;not null"`
	Email     string      `gorm:"type:char(100);unique;uniqueIndex;not null" auth:"username"`
	Image     null.String `gorm:"type:char(100);default:null" json:"-"`
	Password  string      `gorm:"type:char(60);not null" auth:"password" json:"-"`
}

// BeforeCreate hook executed before a User record is inserted in the database.
// Ensures the password is encrypted using bcrypt, with the cost defined by the
// config entry "app.bcryptCost".
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.bcryptPassword(tx)
}

// BeforeUpdate hook executed before a User record is updated in the database.
// Ensures the password is encrypted using bcrypt, with the cost defined by the
// config entry "app.bcryptCost".
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Password") {
		return u.bcryptPassword(tx)
	}

	return nil
}

func (u *User) bcryptPassword(tx *gorm.DB) error {
	var newPass string
	switch u := tx.Statement.Dest.(type) {
	case map[string]interface{}:
		newPass = u["password"].(string)
	case *User:
		newPass = u.Password
	case []*User:
		newPass = u[tx.Statement.CurDestIndex].Password
	}

	b, err := bcrypt.GenerateFromPassword([]byte(newPass), config.GetInt("app.bcryptCost"))
	if err != nil {
		return err
	}
	tx.Statement.SetColumn("password", b)

	return nil
}

// RateLimiterFunc returns rate limiting configuration
// Anonymous users have a quota of 50 requests per minute while
// authenticated users are limited to 500 requests per minute.
func RateLimiterFunc(request *goyave.Request) ratelimiter.Config {
	var id interface{} = nil
	quota := 50
	if request.User != nil {
		id = request.User.(*User).ID
		quota = 500
	}
	return ratelimiter.Config{
		ClientID:      id,
		RequestQuota:  quota,
		QuotaDuration: time.Minute,
	}
}

// UserGenerator generator function for the User model.
// Generate users using the following:
//  database.NewFactory(model.UserGenerator).Generate(5)
func UserGenerator() interface{} {
	user := &User{}
	user.Username = faker.Name()

	b, _ := bcrypt.GenerateFromPassword([]byte(faker.Password()), config.GetInt("app.bcryptCost"))
	user.Password = string(b)

	faker.SetGenerateUniqueValues(true)
	user.Email = faker.Email()
	faker.SetGenerateUniqueValues(false)

	user.Username = faker.Name()
	return user
}
