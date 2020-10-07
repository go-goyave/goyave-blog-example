package model

import (
	"github.com/System-Glitch/goyave/v3/database"
	"github.com/bxcodec/faker/v3"
	"gorm.io/gorm"
)

// A model is a structure reflecting a database table structure. An instance of a model
// is a single database record. Each model is defined in its own file inside the database/models directory.
// Models are usually just normal Golang structs, basic Go types, or pointers of them.
// "sql.Scanner" and "driver.Valuer" interfaces are also supported.

// Learn more here: https://system-glitch.github.io/goyave/guide/basics/database.html#models

func init() {
	// All models should be registered in an "init()" function inside their model file.
	database.RegisterModel(&User{})
}

// User represents a user.
type User struct {
	gorm.Model
	Name  string `gorm:"type:char(100)"`
	Email string `gorm:"type:char(100);unique_index"`
}

// You may need to test features interacting with your database.
// Goyave provides a handy way to generate and save records in your database: factories.
// Factories need a generator function. These functions generate a single random record.
//
// "database.Generator" is an alias for "func() interface{}"
//
// Learn more here: https://system-glitch.github.io/goyave/guide/advanced/testing.html#database-testing

// UserGenerator generator function for the User model.
// Generate users using the following:
//  database.NewFactory(model.UserGenerator).Generate(5)
func UserGenerator() interface{} {
	user := &User{}
	user.Name = faker.Name()

	faker.SetGenerateUniqueValues(true)
	user.Email = faker.Email()
	faker.SetGenerateUniqueValues(false)
	return user
}
