package test

import (
	"net/http"
	"testing"

	"github.com/System-Glitch/goyave/v3"

	"github.com/System-Glitch/goyave-blog-example/database/model"
	"github.com/System-Glitch/goyave-blog-example/http/controller/article"
	"github.com/System-Glitch/goyave-blog-example/http/route"
	_ "github.com/System-Glitch/goyave-blog-example/http/validation"
	"github.com/System-Glitch/goyave/v3/database"
	_ "github.com/System-Glitch/goyave/v3/database/dialect/mysql"
)

type ArticleTestSuite struct {
	goyave.TestSuite
	userID uint
}

type PaginationExpectation struct {
	MaxPage       float64
	Total         float64
	PageSize      float64
	CurrentPage   float64
	RecordsLength float64
}

func (suite *ArticleTestSuite) SetupTest() {
	suite.ClearDatabase()
	factory := database.NewFactory(model.UserGenerator)
	override := &model.User{
		Username: "jack",
		Email:    "jack@example.org",
	}
	suite.userID = factory.Override(override).Save(1).([]*model.User)[0].ID
}

func (suite *ArticleTestSuite) Expect(resp *http.Response, expectation PaginationExpectation) {
	json := map[string]interface{}{}
	err := suite.GetJSONBody(resp, &json)
	suite.Nil(err)
	if err == nil {
		suite.Equal(expectation.MaxPage, json["MaxPage"])
		suite.Equal(expectation.Total, json["Total"])
		suite.Equal(expectation.PageSize, json["PageSize"])
		suite.Equal(expectation.CurrentPage, json["CurrentPage"])

		records := json["Records"].([]interface{})
		suite.Len(records, int(expectation.RecordsLength))
	}
}

func (suite *ArticleTestSuite) TestIndex() {
	suite.RunServer(route.Register, func() {
		factory := database.NewFactory(model.ArticleGenerator)
		override := &model.Article{
			AuthorID: suite.userID,
		}
		factory.Override(override).Save(article.DefaultPageSize + 1)

		resp, err := suite.Get("/article", nil)
		suite.Nil(err)
		if err == nil {
			defer resp.Body.Close()
			suite.Expect(resp, PaginationExpectation{
				MaxPage:       2,
				Total:         article.DefaultPageSize + 1,
				PageSize:      article.DefaultPageSize,
				CurrentPage:   1,
				RecordsLength: article.DefaultPageSize,
			})
		}

		resp, err = suite.Get("/article?page=2", nil)
		suite.Nil(err)
		if err == nil {
			defer resp.Body.Close()
			suite.Expect(resp, PaginationExpectation{
				MaxPage:       2,
				Total:         article.DefaultPageSize + 1,
				PageSize:      article.DefaultPageSize,
				CurrentPage:   2,
				RecordsLength: 1,
			})
		}

		resp, err = suite.Get("/article?pageSize=15&page=1", nil)
		suite.Nil(err)
		if err == nil {
			defer resp.Body.Close()
			suite.Expect(resp, PaginationExpectation{
				MaxPage:       1,
				Total:         article.DefaultPageSize + 1,
				PageSize:      15,
				CurrentPage:   1,
				RecordsLength: 11,
			})
		}
	})
}

func (suite *ArticleTestSuite) TestIndexSearch() {
	suite.RunServer(route.Register, func() {
		factory := database.NewFactory(model.ArticleGenerator)
		override := &model.Article{
			AuthorID: suite.userID,
		}
		factory.Override(override).Save(article.DefaultPageSize + 1)

		override = &model.Article{
			AuthorID: suite.userID,
			Title:    "A very interesting article",
		}
		factory.Override(override).Save(1)

		resp, err := suite.Get("/article?search=interesting", nil)
		suite.Nil(err)
		if err == nil {
			defer resp.Body.Close()
			suite.Expect(resp, PaginationExpectation{
				MaxPage:       1,
				Total:         1,
				PageSize:      article.DefaultPageSize,
				CurrentPage:   1,
				RecordsLength: 1,
			})
		}
	})
}

func TestArticleSuite(t *testing.T) {
	goyave.RunTest(t, new(ArticleTestSuite))
}
