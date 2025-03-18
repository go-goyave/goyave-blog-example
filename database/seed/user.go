package seed

import (
	"math/rand"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/go-goyave/goyave-blog-example/database/model"
	"github.com/samber/lo"
	"goyave.dev/goyave/v5/util/typeutil"
)

func UserGenerator() *model.User {
	user := &model.User{}
	user.Username = typeutil.NewUndefined(faker.Name())
	user.Email = typeutil.NewUndefined(faker.Email(options.WithGenerateUniqueValues(true)))
	user.Password = "$2a$10$TllZ98eJjoknEcE25qR3J.kaLGlOTztt/2SMgbZiTZq5L1O35v76a" // p4ssW0rd_
	user.Articles = lo.Times(rand.Intn(5), func(_ int) *model.Article {
		return ArticleGenerator()
	})
	return user
}
