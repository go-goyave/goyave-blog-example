package model

import (
	"fmt"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
	"goyave.dev/goyave/v4/database"
)

func init() {
	database.RegisterModel(&Article{})
	slug.MaxLength = 80
}

// Article represents an article posted by a user.
type Article struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string `gorm:"type:char(200);not null"`
	Contents  string `gorm:"type:text;not null"`
	Slug      string `gorm:"type:char(80);not null;unique;uniqueIndex"`
	AuthorID  uint   `json:"-"`
	Author    *User  `gorm:"constraint:OnDelete:CASCADE;" json:",omitempty"`
}

// BeforeCreate hook executed before a Article record is inserted in the database.
// Ensures the slug is up to date.
func (a *Article) BeforeCreate(tx *gorm.DB) error {
	return a.slugify(tx)
}

// BeforeUpdate hook executed before an Article record is updated in the database.
// Ensures the slug is up to date.
func (a *Article) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Title") {
		return a.slugify(tx)
	}

	return nil
}

func (a *Article) slugify(tx *gorm.DB) error {
	var newSlug string
	switch a := tx.Statement.Dest.(type) {
	case map[string]interface{}:
		newSlug = generateSlug(a["title"].(string))
	case *Article:
		newSlug = generateSlug(a.Title)
	case []*Article:
		newSlug = generateSlug(a[tx.Statement.CurDestIndex].Title)
	}

	tx.Statement.SetColumn("slug", newSlug)

	return nil
}

// generateSlug creates a slug from the given title. This functions ensures
// the slug is unique by prepending an incremented counter if the slug already
// exists.
func generateSlug(title string) string {
	actualTitle := title
	increment := 0
	s := ""
	count := int64(1)
	for count > 0 {
		s = slug.Make(actualTitle)
		if err := database.Conn().Model(&Article{}).Where("slug = ?", s).Count(&count).Error; err != nil {
			panic(err)
		}
		increment++
		actualTitle = fmt.Sprintf("%d %s", increment, title)
	}
	return s
}

// ArticleGenerator generator function for the Article model.
//
// Be careful, this generator doesn't set the AuthorID!
//
// Generate articles using the following:
//  database.NewFactory(model.ArticleGenerator).Generate(5)
func ArticleGenerator() interface{} {
	article := &Article{}

	faker.SetGenerateUniqueValues(true)
	article.Title = faker.Sentence()
	faker.SetGenerateUniqueValues(false)

	article.Contents = faker.Paragraph()
	article.Slug = slug.Make(article.Title)

	return article
}
