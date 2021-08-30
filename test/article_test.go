package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"goyave.dev/goyave/v4"

	"github.com/go-goyave/goyave-blog-example/database/model"
	"github.com/go-goyave/goyave-blog-example/http/controller/article"
	"github.com/go-goyave/goyave-blog-example/http/route"
	_ "github.com/go-goyave/goyave-blog-example/http/validation"
	"goyave.dev/goyave/v4/auth"
	"goyave.dev/goyave/v4/database"
	_ "goyave.dev/goyave/v4/database/dialect/mysql"
)

type ArticleTestSuite struct {
	goyave.TestSuite
	userID uint
}

type paginationExpectation struct {
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

func (suite *ArticleTestSuite) expect(resp *http.Response, expectation paginationExpectation) {
	json := map[string]interface{}{}
	err := suite.GetJSONBody(resp, &json)
	suite.Nil(err)
	if err == nil {
		suite.Equal(expectation.MaxPage, json["maxPage"])
		suite.Equal(expectation.Total, json["total"])
		suite.Equal(expectation.PageSize, json["pageSize"])
		suite.Equal(expectation.CurrentPage, json["currentPage"])

		records := json["records"].([]interface{})
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
			suite.expect(resp, paginationExpectation{
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
			suite.expect(resp, paginationExpectation{
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
			suite.expect(resp, paginationExpectation{
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
			suite.expect(resp, paginationExpectation{
				MaxPage:       1,
				Total:         1,
				PageSize:      article.DefaultPageSize,
				CurrentPage:   1,
				RecordsLength: 1,
			})
		}
	})
}

func (suite *ArticleTestSuite) TestShow() {
	suite.RunServer(route.Register, func() {
		factory := database.NewFactory(model.ArticleGenerator)
		override := &model.Article{
			AuthorID: suite.userID,
			Title:    "A very interesting article",
		}
		factory.Override(override).Save(1)

		resp, err := suite.Get("/article/a-very-interesting-article", nil)
		suite.Nil(err)
		if err == nil {
			defer resp.Body.Close()
			json := map[string]interface{}{}
			err := suite.GetJSONBody(resp, &json)
			suite.Nil(err)
			if err == nil {
				suite.Equal(override.Title, json["Title"])
			}
		}
	})
}

func (suite *ArticleTestSuite) TestStore() {
	suite.RunServer(route.Register, func() {
		token, err := auth.GenerateToken("jack@example.org")
		if err != nil {
			suite.Error(err)
			return
		}

		headers := map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + token,
		}
		request := map[string]interface{}{
			"title":    "A very interesting article",
			"contents": "lorem ipsum sit dolor amet",
		}
		body, _ := json.Marshal(request)
		resp, err := suite.Post("/article", headers, bytes.NewReader(body))
		suite.Nil(err)
		suite.NotNil(resp)
		if resp != nil {
			defer resp.Body.Close()
			suite.Equal(http.StatusCreated, resp.StatusCode)
			json := map[string]interface{}{}
			err := suite.GetJSONBody(resp, &json)
			suite.Nil(err)
			if err == nil {
				suite.Contains(json, "id")
				suite.Contains(json, "slug")
				suite.Equal("a-very-interesting-article", json["slug"])
			}

			count := int64(0)
			res := database.Conn().
				Model(&model.Article{}).
				Where("slug = ?", "a-very-interesting-article").
				Where("author_id = ?", suite.userID).
				Count(&count)
			if err := res.Error; err != nil {
				suite.Error(err)
			}
			suite.Equal(int64(1), count)
		}
	})
}

func (suite *ArticleTestSuite) TestUpdateByID() {
	suite.RunServer(route.Register, func() {
		factory := database.NewFactory(model.ArticleGenerator)
		override := &model.Article{
			AuthorID: suite.userID,
			Title:    "A very interesting article",
		}
		article := factory.Override(override).Save(1).([]*model.Article)[0]

		suite.testUpdate(fmt.Sprintf("/article/%d", article.ID))
	})
}

func (suite *ArticleTestSuite) TestUpdateBySlug() {
	suite.RunServer(route.Register, func() {
		factory := database.NewFactory(model.ArticleGenerator)
		override := &model.Article{
			AuthorID: suite.userID,
			Title:    "A very interesting article",
		}
		article := factory.Override(override).Save(1).([]*model.Article)[0]

		suite.testUpdate(fmt.Sprintf("/article/%s", article.Slug))
	})
}

func (suite *ArticleTestSuite) testUpdate(url string) {
	token, err := auth.GenerateToken("jack@example.org")
	if err != nil {
		suite.Error(err)
		return
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}
	request := map[string]interface{}{
		"title": "A boring article",
	}
	body, _ := json.Marshal(request)
	resp, err := suite.Patch(url, headers, bytes.NewReader(body))
	suite.Nil(err)
	suite.NotNil(resp)
	if resp != nil {
		defer resp.Body.Close()
		suite.Equal(http.StatusNoContent, resp.StatusCode)

		count := int64(0)
		res := database.Conn().
			Model(&model.Article{}).
			Where("slug = ?", "a-boring-article").
			Count(&count)
		if err := res.Error; err != nil {
			suite.Error(err)
		}
		suite.Equal(int64(1), count)
	}
}

func (suite *ArticleTestSuite) TestDestroyByID() {
	suite.RunServer(route.Register, func() {
		factory := database.NewFactory(model.ArticleGenerator)
		override := &model.Article{
			AuthorID: suite.userID,
			Title:    "A very interesting article",
		}
		article := factory.Override(override).Save(1).([]*model.Article)[0]

		suite.testDestroy(fmt.Sprintf("/article/%d", article.ID), article.ID)
	})
}

func (suite *ArticleTestSuite) TestDestroyBySlug() {
	suite.RunServer(route.Register, func() {
		factory := database.NewFactory(model.ArticleGenerator)
		override := &model.Article{
			AuthorID: suite.userID,
			Title:    "A very interesting article",
		}
		article := factory.Override(override).Save(1).([]*model.Article)[0]

		suite.testDestroy(fmt.Sprintf("/article/%s", article.Slug), article.ID)
	})
}

func (suite *ArticleTestSuite) testDestroy(url string, articleID uint) {
	token, err := auth.GenerateToken("jack@example.org")
	if err != nil {
		suite.Error(err)
		return
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}
	resp, err := suite.Delete(url, headers, nil)
	suite.Nil(err)
	suite.NotNil(resp)
	if resp != nil {
		defer resp.Body.Close()
		suite.Equal(http.StatusNoContent, resp.StatusCode)

		count := int64(0)
		res := database.Conn().
			Model(&model.Article{}).
			Where("id = ?", articleID).
			Count(&count)
		if err := res.Error; err != nil {
			suite.Error(err)
		}
		suite.Equal(int64(0), count)
	}
}

func TestArticleSuite(t *testing.T) {
	goyave.RunTest(t, new(ArticleTestSuite))
}
