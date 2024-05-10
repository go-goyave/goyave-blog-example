package seed

import (
	"gorm.io/gorm"
	"goyave.dev/goyave/v5/database"
)

func Seed(db *gorm.DB) {
	userFactory := database.NewFactory(UserGenerator)
	userFactory.Save(db, 10)
}
