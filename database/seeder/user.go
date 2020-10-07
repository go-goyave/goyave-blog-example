package seeder

import (
	"goyave-blog-example/database/model"

	"github.com/System-Glitch/goyave/v3/database"
	"github.com/bxcodec/faker/v3"
)

// Seeders are functions which create a number of random records in the database
// in order to create a full and realistic test environment.
//
// Each seeder should have its own file.
// A seeder's responsibilities are limited to a single table or model.
// For example, the "seeder.User" should only seed the "users" table.
// Moreover, seeders should have the same name as the model they are using.
//
// Learn more here: https://system-glitch.github.io/goyave/guide/advanced/testing.html#seeders

// User seeder for users. Generate and save 10 users in the database.
func User() {
	database.NewFactory(model.UserGenerator).Save(10)

	// As user generator makes unique emails,
	// forget generated unique emails.
	// See https://github.com/bxcodec/faker/blob/master/SingleFakeData.md#unique-values
	faker.ResetUnique()
}
