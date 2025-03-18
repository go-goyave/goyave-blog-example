package seed

import (
	"github.com/go-faker/faker/v4"
	"github.com/go-goyave/goyave-blog-example/database/model"
	"github.com/go-goyave/goyave-blog-example/service/article"
	"github.com/samber/lo"
	"goyave.dev/goyave/v5/util/typeutil"
)

func ArticleGenerator() *model.Article {
	a := &model.Article{}
	a.Title = typeutil.NewUndefined(faker.Sentence())
	a.Contents = typeutil.NewUndefined(faker.Paragraph())
	a.Slug = typeutil.NewUndefined(lo.Must(article.NewService(nil, nil).GenerateSlug(a.Title.Val)))
	return a
}
