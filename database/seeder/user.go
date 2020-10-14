package seeder

import (
	"github.com/System-Glitch/goyave-blog-example/database/model"

	"github.com/System-Glitch/goyave/v3/database"
	"github.com/bxcodec/faker/v3"
)

// User seeder for users. Generate and save 10 users in the database.
func User() {
	database.NewFactory(model.UserGenerator).Save(10)

	// As user generator makes unique emails,
	// forget generated unique emails.
	// See https://github.com/bxcodec/faker/blob/master/SingleFakeData.md#unique-values
	faker.ResetUnique()
}
