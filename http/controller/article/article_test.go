package article

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-goyave/goyave-blog-example/database/seed"
	"github.com/go-goyave/goyave-blog-example/dto"
	"github.com/go-goyave/goyave-blog-example/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"goyave.dev/filter"
	"goyave.dev/goyave/v5"
	"goyave.dev/goyave/v5/auth"
	"goyave.dev/goyave/v5/database"
	"goyave.dev/goyave/v5/middleware/parse"
	"goyave.dev/goyave/v5/util/testutil"
	"goyave.dev/goyave/v5/util/typeutil"
)

type updateArticleDTO struct {
	Title    string `json:"title"`
	Contents string `json:"contents"`
}

type serviceMock struct {
	paginator *database.PaginatorDTO[*dto.Article]
	article   *dto.Article
	err       error

	createCallback func(*dto.CreateArticle)
	updateCallback func(*dto.UpdateArticle)

	isOwner bool
}

func (s *serviceMock) Index(_ context.Context, _ *filter.Request) (*database.PaginatorDTO[*dto.Article], error) {
	return s.paginator, s.err
}

func (s *serviceMock) GetBySlug(_ context.Context, slug string) (*dto.Article, error) {
	if s.article.Slug == slug {
		return s.article, s.err
	}
	return nil, gorm.ErrRecordNotFound
}

func (s *serviceMock) Create(_ context.Context, createDTO *dto.CreateArticle) error {
	s.createCallback(createDTO)
	return s.err
}

func (s *serviceMock) Update(_ context.Context, _ uint, updateDTO *dto.UpdateArticle) error {
	s.updateCallback(updateDTO)
	return s.err
}

func (s *serviceMock) Delete(_ context.Context, _ uint) error {
	return s.err
}

func (s *serviceMock) IsOwner(_ context.Context, _ uint, _ uint) (bool, error) {
	return s.isOwner, nil
}

func (s *serviceMock) Name() string {
	return service.Article
}

const mockAuthUserMeta = "mock:authuser"

type mockAuthMiddleware struct {
	goyave.Component
}

func (m *mockAuthMiddleware) Handle(next goyave.Handler) goyave.Handler {
	return func(response *goyave.Response, request *goyave.Request) {
		request.User, _ = request.Route.LookupMeta(mockAuthUserMeta)
		requireAuth, _ := request.Route.LookupMeta(auth.MetaAuth)
		if requireAuth.(bool) && request.User == nil {
			response.Status(http.StatusUnauthorized)
			return
		}
		next(response, request)
	}
}

func generatePaginator() *database.PaginatorDTO[*dto.Article] {
	records := database.NewFactory(seed.ArticleGenerator).Generate(3)
	return &database.PaginatorDTO[*dto.Article]{
		Records:     typeutil.MustConvert[[]*dto.Article](records),
		MaxPage:     1,
		Total:       3,
		PageSize:    10,
		CurrentPage: 1,
	}
}

func setupArticleTest(t *testing.T, service *serviceMock) *testutil.TestServer {
	server := testutil.NewTestServer(t, "config.test.json")
	server.RegisterService(service)
	server.RegisterRoutes(func(_ *goyave.Server, r *goyave.Router) {
		r.GlobalMiddleware(&parse.Middleware{})
		r.Controller(NewController())
	})
	return server
}

func TestArticle(t *testing.T) {
	t.Run("Index", func(t *testing.T) {
		service := &serviceMock{
			paginator: generatePaginator(),
		}
		server := setupArticleTest(t, service)
		request := httptest.NewRequest(http.MethodGet, "/articles", nil)
		response := server.TestRequest(request)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		paginator, err := testutil.ReadJSONBody[*database.PaginatorDTO[*dto.Article]](response.Body)
		assert.NoError(t, err)
		assert.NoError(t, response.Body.Close())

		assert.Equal(t, service.paginator, paginator)

		t.Run("error", func(t *testing.T) {
			service.err = fmt.Errorf("test error")
			request := httptest.NewRequest(http.MethodGet, "/articles", nil)
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
			assert.NoError(t, response.Body.Close())
		})
	})

	t.Run("Show", func(t *testing.T) {
		service := &serviceMock{
			article: typeutil.MustConvert[*dto.Article](seed.ArticleGenerator()),
		}
		server := setupArticleTest(t, service)
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/articles/%s", service.article.Slug), nil)
		response := server.TestRequest(request)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		paginator, err := testutil.ReadJSONBody[*dto.Article](response.Body)
		assert.NoError(t, err)
		assert.NoError(t, response.Body.Close())

		assert.Equal(t, service.article, paginator)

		t.Run("not_found", func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/articles/incorrect-slug", nil)
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusNotFound, response.StatusCode)
			assert.NoError(t, response.Body.Close())
		})
	})

	t.Run("Create", func(t *testing.T) {
		service := &serviceMock{}
		server := setupArticleTest(t, service)
		user := &dto.InternalUser{
			User: dto.User{ID: 1},
		}
		server.Router().GlobalMiddleware(&mockAuthMiddleware{}).SetMeta(mockAuthUserMeta, user)

		requestBody := &dto.CreateArticle{
			Title:    "article title",
			Contents: "article contents",
		}

		request := httptest.NewRequest(http.MethodPost, "/articles", testutil.ToJSON(requestBody))
		request.Header.Set("Content-Type", "application/json")

		service.createCallback = func(createDTO *dto.CreateArticle) {
			expected := typeutil.Copy(&dto.CreateArticle{AuthorID: user.ID}, requestBody)
			assert.Equal(t, expected, createDTO)
		}

		response := server.TestRequest(request)
		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.NoError(t, response.Body.Close())

		t.Run("error", func(t *testing.T) {
			service.err = fmt.Errorf("test error")
			request := httptest.NewRequest(http.MethodPost, "/articles", testutil.ToJSON(requestBody))
			request.Header.Set("Content-Type", "application/json")
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
			assert.NoError(t, response.Body.Close())
		})

		t.Run("require_auth", func(t *testing.T) {
			server.Router().RemoveMeta(mockAuthUserMeta)
			request := httptest.NewRequest(http.MethodPost, "/articles", testutil.ToJSON(requestBody))
			request.Header.Set("Content-Type", "application/json")
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
			assert.NoError(t, response.Body.Close())
		})
	})

	t.Run("Update", func(t *testing.T) {
		service := &serviceMock{}
		server := setupArticleTest(t, service)
		user := &dto.InternalUser{
			User: dto.User{ID: 1},
		}
		server.Router().GlobalMiddleware(&mockAuthMiddleware{}).SetMeta(mockAuthUserMeta, user)

		requestBody := &updateArticleDTO{
			Title:    "article title",
			Contents: "article contents",
		}

		request := httptest.NewRequest(http.MethodPatch, "/articles/1", testutil.ToJSON(requestBody))
		request.Header.Set("Content-Type", "application/json")

		service.isOwner = true
		service.updateCallback = func(updateDTO *dto.UpdateArticle) {
			assert.Equal(t, requestBody.Title, updateDTO.Title.Val)
			assert.Equal(t, requestBody.Contents, updateDTO.Contents.Val)
		}

		response := server.TestRequest(request)
		assert.Equal(t, http.StatusNoContent, response.StatusCode)
		assert.NoError(t, response.Body.Close())

		t.Run("invalid_id", func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPatch, "/articles/999999999999999999999999999999999999", testutil.ToJSON(requestBody))
			request.Header.Set("Content-Type", "application/json")
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusNotFound, response.StatusCode)
			assert.NoError(t, response.Body.Close())
		})

		t.Run("not_found", func(t *testing.T) {
			service.err = gorm.ErrRecordNotFound
			request := httptest.NewRequest(http.MethodPatch, "/articles/1", testutil.ToJSON(requestBody))
			request.Header.Set("Content-Type", "application/json")
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusNotFound, response.StatusCode)
			assert.NoError(t, response.Body.Close())
		})

		t.Run("error", func(t *testing.T) {
			service.err = fmt.Errorf("test error")
			request := httptest.NewRequest(http.MethodPatch, "/articles/1", testutil.ToJSON(requestBody))
			request.Header.Set("Content-Type", "application/json")
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
			assert.NoError(t, response.Body.Close())
		})

		t.Run("not_owner", func(t *testing.T) {
			service.err = nil
			service.isOwner = false
			request := httptest.NewRequest(http.MethodPatch, "/articles/1", testutil.ToJSON(requestBody))
			request.Header.Set("Content-Type", "application/json")
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusForbidden, response.StatusCode)
			assert.NoError(t, response.Body.Close())
		})

		t.Run("require_auth", func(t *testing.T) {
			server.Router().RemoveMeta(mockAuthUserMeta)
			request := httptest.NewRequest(http.MethodPatch, "/articles/1", testutil.ToJSON(requestBody))
			request.Header.Set("Content-Type", "application/json")
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
			assert.NoError(t, response.Body.Close())
		})
	})

	t.Run("Delete", func(t *testing.T) {
		service := &serviceMock{}
		server := setupArticleTest(t, service)
		user := &dto.InternalUser{
			User: dto.User{ID: 1},
		}
		server.Router().GlobalMiddleware(&mockAuthMiddleware{}).SetMeta(mockAuthUserMeta, user)
		service.isOwner = true
		request := httptest.NewRequest(http.MethodDelete, "/articles/1", nil)
		response := server.TestRequest(request)
		assert.Equal(t, http.StatusNoContent, response.StatusCode)
		assert.NoError(t, response.Body.Close())

		t.Run("invalid_id", func(t *testing.T) {
			request := httptest.NewRequest(http.MethodDelete, "/articles/999999999999999999999999999999999999", nil)
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusNotFound, response.StatusCode)
			assert.NoError(t, response.Body.Close())
		})

		t.Run("not_found", func(t *testing.T) {
			service.err = gorm.ErrRecordNotFound
			request := httptest.NewRequest(http.MethodDelete, "/articles/2", nil)
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusNotFound, response.StatusCode)
			assert.NoError(t, response.Body.Close())
		})

		t.Run("not_owner", func(t *testing.T) {
			service.err = nil
			service.isOwner = false
			request := httptest.NewRequest(http.MethodDelete, "/articles/1", nil)
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusForbidden, response.StatusCode)
			assert.NoError(t, response.Body.Close())
		})

		t.Run("require_auth", func(t *testing.T) {
			server.Router().RemoveMeta(mockAuthUserMeta)
			request := httptest.NewRequest(http.MethodDelete, "/articles/1", nil)
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
			assert.NoError(t, response.Body.Close())
		})
	})
}
