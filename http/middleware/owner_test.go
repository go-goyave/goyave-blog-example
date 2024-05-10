package middleware

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"goyave.dev/goyave/v5"
	"goyave.dev/goyave/v5/slog"
	"goyave.dev/goyave/v5/util/testutil"

	"github.com/go-goyave/goyave-blog-example/dto"
	"github.com/stretchr/testify/assert"
)

type ownerServiceMock struct {
	err     error
	isOwner bool
}

func (s ownerServiceMock) IsOwner(_ context.Context, _, _ uint) (bool, error) {
	return s.isOwner, s.err
}

func TestOwner(t *testing.T) {

	cases := []struct {
		service    ownerServiceMock
		articleID  string
		desc       string
		wantStatus int
	}{
		{desc: "OK", service: ownerServiceMock{isOwner: true}, articleID: "123", wantStatus: http.StatusOK},
		{desc: "NOK", service: ownerServiceMock{isOwner: false}, articleID: "123", wantStatus: http.StatusForbidden},
		{desc: "record_not_found", service: ownerServiceMock{}, articleID: "NaN", wantStatus: http.StatusNotFound},
		{desc: "db_error", service: ownerServiceMock{err: fmt.Errorf("test db error")}, articleID: "123", wantStatus: http.StatusInternalServerError},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			middleware := NewOwner("articleID", c.service)
			server := testutil.NewTestServer(t, "config.test.json")
			server.Logger = slog.New(slog.NewHandler(false, io.Discard))
			request := server.NewTestRequest(http.MethodGet, "/article/"+c.articleID, nil)
			request.RouteParams = map[string]string{"articleID": c.articleID}
			request.User = &dto.InternalUser{User: dto.User{ID: 1}}
			resp := server.TestMiddleware(middleware, request, func(response *goyave.Response, _ *goyave.Request) {
				response.Status(http.StatusOK)
			})
			assert.NoError(t, resp.Body.Close())
			assert.Equal(t, c.wantStatus, resp.StatusCode)
		})
	}
}
