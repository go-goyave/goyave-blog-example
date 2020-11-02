package dbutil

import (
	"strings"

	"github.com/System-Glitch/goyave-blog-example/database/model"
	"github.com/System-Glitch/goyave-blog-example/database/seeder"
	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/config"
	"github.com/System-Glitch/goyave/v3/database"
	"gorm.io/gorm"
)

// Paginate create a tx scope for pagination.
//  conn.Scopes(database.Paginate(r)).Find(&users)
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// EscapeLike escape "%" and "_" characters in the given string
// for use in "LIKE" clauses.
func EscapeLike(str string) string {
	escapeChars := []string{"%", "_"}
	for _, v := range escapeChars {
		str = strings.ReplaceAll(str, v, "\\"+v)
	}
	return str
}

// RunSeeders run seeders if the user table is empty.
// Only triggers if the environment is "localhost".
func RunSeeders() {
	if config.GetString("app.environment") == "localhost" {
		count := int64(0)
		if err := database.Conn().Model(&model.User{}).Count(&count).Error; err != nil {
			panic(err)
		}

		if count <= 0 {
			goyave.Logger.Println("Running seeders...")
			seeder.User()
			seeder.Article()
		}
	}
}
