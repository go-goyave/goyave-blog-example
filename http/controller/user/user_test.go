package user

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"regexp"
	"testing"
	"time"

	"github.com/go-goyave/goyave-blog-example/dto"
	"github.com/go-goyave/goyave-blog-example/service"
	"github.com/go-goyave/goyave-blog-example/service/storage"
	"github.com/guregu/null/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"goyave.dev/goyave/v5"
	"goyave.dev/goyave/v5/auth"
	"goyave.dev/goyave/v5/middleware/parse"
	"goyave.dev/goyave/v5/util/fsutil/osfs"
	"goyave.dev/goyave/v5/util/testutil"
	"goyave.dev/goyave/v5/util/typeutil"

	"github.com/DATA-DOG/go-sqlmock"
)

type upsertDTO struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type serviceMock struct {
	user             *dto.InternalUser
	registerCallback func(*dto.RegisterUser)
	updateCallback   func(*dto.UpdateUser)
	err              error
}

func (s *serviceMock) UniqueScope() func(db *gorm.DB, val any) *gorm.DB {
	return func(db *gorm.DB, val any) *gorm.DB {
		return db.Table("users").Where("email", val)
	}
}

func (s *serviceMock) GetByID(_ context.Context, _ uint) (*dto.InternalUser, error) {
	return s.user, s.err
}

func (s *serviceMock) Register(_ context.Context, registerDTO *dto.RegisterUser) error {
	s.registerCallback(registerDTO)
	return s.err
}

func (s *serviceMock) Update(_ context.Context, _ uint, updateDTO *dto.UpdateUser) error {
	s.updateCallback(updateDTO)
	return s.err
}

func (s *serviceMock) Name() string {
	return service.User
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

func setupUserTest(t *testing.T, service *serviceMock) *testutil.TestServer {
	server := testutil.NewTestServer(t, "config.test.json")
	server.RegisterService(service)

	rootDir := testutil.FindRootDirectory()

	imgFS := osfs.New(path.Join(rootDir, "resources/img"))
	storageService := storage.NewService(imgFS, imgFS)
	server.RegisterService(storageService)
	server.RegisterRoutes(func(_ *goyave.Server, r *goyave.Router) {
		r.GlobalMiddleware(&parse.Middleware{})
		r.Controller(NewController())
	})
	return server
}

func setupMock(t *testing.T, server *testutil.TestServer) sqlmock.Sqlmock {
	server.Config().Set("database.config.prepareStmt", false)
	server.Config().Set("database.connection", "mock")
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	dialector := postgres.New(postgres.Config{
		DSN:                  "mock_db",
		DriverName:           "postgres",
		Conn:                 mockDB,
		PreferSimpleProtocol: true,
	})
	require.NoError(t, server.ReplaceDB(dialector))
	return mock
}

func TestUser(t *testing.T) {
	t.Run("ShowProfile", func(t *testing.T) {
		server := setupUserTest(t, &serviceMock{})
		user := &dto.InternalUser{
			User: dto.User{
				ID:        1,
				CreatedAt: time.Now().Round(0).UTC(),
				Username:  "johndoe",
				Email:     "johndoe@example.org",
			},
			Avatar: null.NewString("img.jpeg", true),
		}
		server.Router().GlobalMiddleware(&mockAuthMiddleware{}).SetMeta(mockAuthUserMeta, user)

		request := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
		response := server.TestRequest(request)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		profile, err := testutil.ReadJSONBody[*dto.InternalUser](response.Body)
		assert.NoError(t, err)
		assert.NoError(t, response.Body.Close())
		assert.Equal(t, &dto.InternalUser{User: user.User}, profile)
	})

	t.Run("ShowAvatar", func(t *testing.T) {
		service := &serviceMock{}
		server := setupUserTest(t, service)
		user := &dto.InternalUser{
			User:   dto.User{ID: 1},
			Avatar: null.NewString("test_profile_picture.jpg", true),
		}
		service.user = user

		imgFile, err := osfs.New(path.Join(testutil.FindRootDirectory(), "resources/img")).Open("test_profile_picture.jpg")
		require.NoError(t, err)
		profilePicture, err := io.ReadAll(imgFile)
		assert.NoError(t, imgFile.Close())
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodGet, "/users/1/avatar", nil)
		response := server.TestRequest(request)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		responseProfilePicture, err := io.ReadAll(response.Body)
		assert.NoError(t, err)
		assert.NoError(t, response.Body.Close())
		assert.Equal(t, profilePicture, responseProfilePicture)

		t.Run("default_avatar", func(t *testing.T) {
			service.user = &dto.InternalUser{User: user.User}
			imgFile, err := osfs.New(path.Join(testutil.FindRootDirectory(), "resources/img")).Open("default_profile_picture.jpg")
			require.NoError(t, err)
			defaultProfilePicture, err := io.ReadAll(imgFile)
			assert.NoError(t, imgFile.Close())
			require.NoError(t, err)

			request := httptest.NewRequest(http.MethodGet, "/users/1/avatar", nil)
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusOK, response.StatusCode)
			responseProfilePicture, err := io.ReadAll(response.Body)
			assert.NoError(t, err)
			assert.NoError(t, response.Body.Close())
			assert.Equal(t, defaultProfilePicture, responseProfilePicture)
		})

		t.Run("invalid_id", func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/users/999999999999999999999999999999999999/avatar", nil)
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusNotFound, response.StatusCode)
			assert.NoError(t, response.Body.Close())
		})

		t.Run("error", func(t *testing.T) {
			service.err = fmt.Errorf("test error")
			request := httptest.NewRequest(http.MethodGet, "/users/1/avatar", nil)
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
			assert.NoError(t, response.Body.Close())
		})

		t.Run("not_found", func(t *testing.T) {
			service.err = gorm.ErrRecordNotFound
			request := httptest.NewRequest(http.MethodGet, "/users/1/avatar", nil)
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusNotFound, response.StatusCode)
			assert.NoError(t, response.Body.Close())
		})
	})

	t.Run("Register", func(t *testing.T) {
		service := &serviceMock{}
		server := setupUserTest(t, service)
		mock := setupMock(t, server)

		requestBody := &upsertDTO{
			Email:    "johndoe@example.org",
			Username: "johndoe",
			Password: "p4ssW0rd_",
		}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM \"users\" WHERE \"email\" = $1")).
			WithArgs(requestBody.Email).
			WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0))
		defer func() {
			assert.NoError(t, mock.ExpectationsWereMet())
		}()

		request := httptest.NewRequest(http.MethodPost, "/users", testutil.ToJSON(requestBody))
		request.Header.Set("Content-Type", "application/json")

		service.registerCallback = func(registerDTO *dto.RegisterUser) {
			expected := typeutil.Copy(&dto.RegisterUser{}, requestBody)
			expected.Password = requestBody.Password
			assert.Equal(t, expected, registerDTO)
		}

		response := server.TestRequest(request)
		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.NoError(t, response.Body.Close())

		t.Run("error", func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM \"users\" WHERE \"email\" = $1")).
				WithArgs(requestBody.Email).
				WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0))
			service.err = fmt.Errorf("test error")
			request := httptest.NewRequest(http.MethodPost, "/users", testutil.ToJSON(requestBody))
			request.Header.Set("Content-Type", "application/json")
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
			assert.NoError(t, response.Body.Close())
		})
	})

	t.Run("Update", func(t *testing.T) {
		service := &serviceMock{}
		server := setupUserTest(t, service)
		user := &dto.InternalUser{
			User: dto.User{
				ID:       1,
				Username: "johndoe",
				Email:    "johndoe@example.org",
			},
		}
		server.Router().GlobalMiddleware(&mockAuthMiddleware{}).SetMeta(mockAuthUserMeta, user)
		mock := setupMock(t, server)

		requestBody := &upsertDTO{
			Email:    "johndoe-updated@example.org",
			Username: "johndoe-updated",
			Password: "new-p4ssW0rd_",
		}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM \"users\" WHERE \"email\" = $1")).
			WithArgs(requestBody.Email).
			WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0))
		defer func() {
			assert.NoError(t, mock.ExpectationsWereMet())
		}()

		request := httptest.NewRequest(http.MethodPatch, "/users", testutil.ToJSON(requestBody))
		request.Header.Set("Content-Type", "application/json")

		service.updateCallback = func(registerDTO *dto.UpdateUser) {
			expected := typeutil.Copy(&dto.UpdateUser{}, requestBody)
			expected.Password = typeutil.NewUndefined(requestBody.Password)
			assert.Equal(t, expected, registerDTO)
		}

		response := server.TestRequest(request)
		assert.Equal(t, http.StatusNoContent, response.StatusCode)
		assert.NoError(t, response.Body.Close())

		t.Run("error", func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM \"users\" WHERE \"email\" = $1")).
				WithArgs(requestBody.Email).
				WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0))
			service.err = fmt.Errorf("test error")
			request := httptest.NewRequest(http.MethodPatch, "/users", testutil.ToJSON(requestBody))
			request.Header.Set("Content-Type", "application/json")
			response := server.TestRequest(request)
			assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
			assert.NoError(t, response.Body.Close())
		})
	})
}
