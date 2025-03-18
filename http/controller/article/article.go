package article

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-goyave/goyave-blog-example/dto"
	"github.com/go-goyave/goyave-blog-example/http/middleware"
	"github.com/go-goyave/goyave-blog-example/service"
	"goyave.dev/filter"
	"goyave.dev/goyave/v5"
	"goyave.dev/goyave/v5/auth"
	"goyave.dev/goyave/v5/database"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Service interface {
	Index(ctx context.Context, request *filter.Request) (*database.PaginatorDTO[*dto.Article], error)
	GetBySlug(ctx context.Context, slug string) (*dto.Article, error)
	Create(ctx context.Context, createDTO *dto.CreateArticle) error
	Update(ctx context.Context, id int64, updateDTO *dto.UpdateArticle) error
	Delete(ctx context.Context, id int64) error
	IsOwner(ctx context.Context, resourceID, ownerID int64) (bool, error)
}

type Controller struct {
	goyave.Component
	ArticleService Service
}

func NewController() *Controller {
	return &Controller{}
}

func (ctrl *Controller) Init(server *goyave.Server) {
	ctrl.Component.Init(server)
	ctrl.ArticleService = server.Service(service.Article).(Service)
}

func (ctrl *Controller) RegisterRoutes(router *goyave.Router) {
	subrouter := router.Subrouter("/articles")
	subrouter.Get("/", ctrl.Index).ValidateQuery(filter.Validation)
	subrouter.Get("/{slug}", ctrl.Show)

	authRouter := subrouter.Group().SetMeta(auth.MetaAuth, true)
	authRouter.Post("/", ctrl.Create).ValidateBody(ctrl.CreateRequest)

	ownedRouter := authRouter.Group()
	ownerMiddleware := middleware.NewOwner("articleID", ctrl.ArticleService)
	ownedRouter.Middleware(ownerMiddleware)
	ownedRouter.Patch("/{articleID:[0-9]+}", ctrl.Update).ValidateBody(ctrl.UpdateRequest)
	ownedRouter.Delete("/{articleID:[0-9]+}", ctrl.Delete)
}

func (ctrl *Controller) Index(response *goyave.Response, request *goyave.Request) {
	paginator, err := ctrl.ArticleService.Index(request.Context(), filter.NewRequest(request.Query))
	if response.WriteDBError(err) {
		return
	}
	response.JSON(http.StatusOK, paginator)
}

func (ctrl *Controller) Show(response *goyave.Response, request *goyave.Request) {
	user, err := ctrl.ArticleService.GetBySlug(request.Context(), request.RouteParams["slug"])
	if response.WriteDBError(err) {
		return
	}
	response.JSON(http.StatusOK, user)
}

func (ctrl *Controller) Create(response *goyave.Response, request *goyave.Request) {
	createDTO := typeutil.MustConvert[*dto.CreateArticle](request.Data)
	createDTO.AuthorID = request.User.(*dto.InternalUser).ID.Val

	err := ctrl.ArticleService.Create(request.Context(), createDTO)
	if err != nil {
		response.Error(err)
		return
	}
	response.Status(http.StatusCreated)
}

func (ctrl *Controller) Update(response *goyave.Response, request *goyave.Request) {
	id, err := strconv.ParseInt(request.RouteParams["articleID"], 10, 64)
	if err != nil {
		response.Status(http.StatusNotFound)
		return
	}

	updateDTO := typeutil.MustConvert[*dto.UpdateArticle](request.Data)

	err = ctrl.ArticleService.Update(request.Context(), id, updateDTO)
	response.WriteDBError(err)
}

func (ctrl *Controller) Delete(response *goyave.Response, request *goyave.Request) {
	id, err := strconv.ParseInt(request.RouteParams["articleID"], 10, 64)
	if err != nil {
		response.Status(http.StatusNotFound)
		return
	}

	err = ctrl.ArticleService.Delete(request.Context(), id)
	response.WriteDBError(err)
}
