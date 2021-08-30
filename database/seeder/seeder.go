package seeder

import (
	"github.com/go-goyave/goyave-blog-example/database/model"
	"goyave.dev/goyave/v4"
	"goyave.dev/goyave/v4/config"
	"goyave.dev/goyave/v4/database"
)

// Run run seeders if the user table is empty.
// Only triggers if the environment is "localhost".
func Run() {
	if config.GetString("app.environment") == "localhost" {
		count := int64(0)
		if err := database.Conn().Model(&model.User{}).Count(&count).Error; err != nil {
			panic(err)
		}

		if count <= 0 {
			goyave.Logger.Println("Running seeders...")
			User()
			Article()
		}
	}
}
