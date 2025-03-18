package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-goyave/goyave-blog-example/dto"
	"goyave.dev/goyave/v5"
)

type OwnerService interface {
	IsOwner(ctx context.Context, resourceID, ownerID int64) (bool, error)
}

type Owner struct {
	goyave.Component

	OwnerService OwnerService

	// RouteParam the name of the route param identifying the requested resource (e.g: "articleID")
	RouteParam string
}

func NewOwner(routeParam string, ownerService OwnerService) *Owner {
	return &Owner{
		RouteParam:   routeParam,
		OwnerService: ownerService,
	}
}

func (m *Owner) Handle(next goyave.Handler) goyave.Handler {
	return func(response *goyave.Response, request *goyave.Request) {
		resourceID, err := strconv.ParseInt(request.RouteParams[m.RouteParam], 10, 64)
		if err != nil {
			response.Status(http.StatusNotFound)
			return
		}

		user := request.User.(*dto.InternalUser)

		isOwner, err := m.OwnerService.IsOwner(request.Context(), resourceID, user.ID.Val)
		if response.WriteDBError(err) {
			return
		}

		if !isOwner {
			response.Status(http.StatusForbidden)
			return
		}

		next(response, request)
	}
}
