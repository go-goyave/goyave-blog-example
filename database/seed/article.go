package seed

import (
	"github.com/go-faker/faker/v4"
	"github.com/go-goyave/goyave-blog-example/database/model"
	"github.com/go-goyave/goyave-blog-example/service/article"
	"github.com/samber/lo"
)

func ArticleGenerator() *model.Article {
	a := &model.Article{}
	a.Title = faker.Sentence()
	a.Contents = faker.Paragraph()
	a.Slug = lo.Must(article.NewService(nil, nil).GenerateSlug(a.Title))
	return a
}
