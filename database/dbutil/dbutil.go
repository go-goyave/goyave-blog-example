package dbutil

import "gorm.io/gorm"

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
