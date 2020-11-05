package dbutil

import (
	"strings"

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
