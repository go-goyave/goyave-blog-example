package article

import (
	"context"
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/go-goyave/goyave-blog-example/database/model"
	"github.com/go-goyave/goyave-blog-example/dto"
	"github.com/go-goyave/goyave-blog-example/service"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"goyave.dev/filter"
	"goyave.dev/goyave/v5/database"
	"goyave.dev/goyave/v5/util/errors"
	"goyave.dev/goyave/v5/util/session"
	"goyave.dev/goyave/v5/util/typeutil"
)

func init() {
	slug.MaxLength = 126
}

type Repository interface {
	Index(ctx context.Context, request *filter.Request) (*database.Paginator[*model.Article], error)
	Create(ctx context.Context, article *model.Article) (*model.Article, error)
	Update(ctx context.Context, article *model.Article) (*model.Article, error)
	GetByID(ctx context.Context, id uint) (*model.Article, error)
	GetBySlug(ctx context.Context, slug string) (*model.Article, error)
	Delete(ctx context.Context, id uint) error
	IsOwner(ctx context.Context, resourceID, ownerID uint) (bool, error)
}

type Service struct {
	Session    session.Session
	Repository Repository
}

func NewService(session session.Session, repository Repository) *Service {
	return &Service{
		Session:    session,
		Repository: repository,
	}
}

func (s *Service) Index(ctx context.Context, request *filter.Request) (*database.PaginatorDTO[*dto.Article], error) {
	paginator, err := s.Repository.Index(ctx, request)
	if err != nil {
		return nil, errors.New(err)
	}
	return typeutil.MustConvert[*database.PaginatorDTO[*dto.Article]](paginator), nil
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*dto.Article, error) {
	user, err := s.Repository.GetBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New(err)
	}
	return typeutil.MustConvert[*dto.Article](user), nil
}

func (s *Service) Create(ctx context.Context, createDTO *dto.CreateArticle) error {
	article := typeutil.Copy(&model.Article{}, createDTO)
	var err error
	article.Slug, err = s.GenerateSlug(article.Title)
	if err != nil {
		return errors.New(err)
	}
	_, err = s.Repository.Create(ctx, article)
	return errors.New(err)
}

func (s *Service) GenerateSlug(title string) (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", errors.New(err)
	}

	shortUID := strings.ToLower(strings.TrimRight(base32.StdEncoding.EncodeToString(uuid[:]), "="))
	return slug.Make(fmt.Sprintf("%s-%s", shortUID, title)), nil
}

func (s *Service) Update(ctx context.Context, id uint, updateDTO *dto.UpdateArticle) error {
	err := s.Session.Transaction(ctx, func(ctx context.Context) error {
		user, err := s.Repository.GetByID(ctx, id)
		if err != nil {
			return errors.New(err)
		}

		user = typeutil.Copy(user, updateDTO)

		_, err = s.Repository.Update(ctx, user)
		return errors.New(err)
	})

	return errors.New(err)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.Repository.Delete(ctx, id)
}

func (s *Service) IsOwner(ctx context.Context, resourceID, ownerID uint) (bool, error) {
	return s.Repository.IsOwner(ctx, resourceID, ownerID)
}

func (s *Service) Name() string {
	return service.Article
}
