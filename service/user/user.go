package user

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/go-goyave/goyave-blog-example/database/model"
	"github.com/go-goyave/goyave-blog-example/dto"
	"github.com/go-goyave/goyave-blog-example/service"
	"github.com/guregu/null/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"goyave.dev/goyave/v5/slog"
	"goyave.dev/goyave/v5/util/errors"
	"goyave.dev/goyave/v5/util/fsutil"
	"goyave.dev/goyave/v5/util/session"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Repository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	UniqueScope() func(db *gorm.DB, val any) *gorm.DB
}

type StorageService interface {
	GetFS() fs.StatFS
	SaveAvatar(file fsutil.File) (string, error)
	Delete(path string) error
}

type Service struct {
	Session        session.Session
	Repository     Repository
	StorageService StorageService
	Logger         *slog.Logger
}

func NewService(session session.Session, logger *slog.Logger, repository Repository, storageService StorageService) *Service {
	return &Service{
		Session:        session,
		Logger:         logger,
		Repository:     repository,
		StorageService: storageService,
	}
}

func (s *Service) UniqueScope() func(db *gorm.DB, val any) *gorm.DB {
	return s.Repository.UniqueScope()
}

func (s *Service) GetByID(ctx context.Context, id uint) (*dto.InternalUser, error) {
	user, err := s.Repository.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New(err)
	}
	return typeutil.MustConvert[*dto.InternalUser](user), nil
}

func (s *Service) FindByUsername(ctx context.Context, username any) (*dto.InternalUser, error) {
	user, err := s.Repository.GetByEmail(ctx, fmt.Sprintf("%v", username))
	if err != nil {
		return nil, errors.New(err)
	}
	return typeutil.MustConvert[*dto.InternalUser](user), nil
}

func (s *Service) Register(ctx context.Context, registerDTO *dto.RegisterUser) error {
	err := s.Session.Transaction(ctx, func(ctx context.Context) error {
		user := typeutil.Copy(&model.User{}, registerDTO)

		b, err := bcrypt.GenerateFromPassword([]byte(registerDTO.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New(err)
		}
		user.Password = string(b)

		if registerDTO.Avatar.IsPresent() && registerDTO.Avatar.Val.Valid {
			filename, err := s.StorageService.SaveAvatar(registerDTO.Avatar.Val.V[0])
			if err != nil {
				return errors.New(err)
			}
			user.Avatar.SetValid(filename)
		}

		_, err = s.Repository.Create(ctx, user)
		if err != nil && user.Avatar.Valid {
			if err := s.StorageService.Delete(user.Avatar.String); err != nil {
				s.Logger.Error(errors.New(err))
			}
		}
		return errors.New(err)
	})
	return errors.New(err)
}

func (s *Service) Update(ctx context.Context, id uint, updateDTO *dto.UpdateUser) error {
	err := s.Session.Transaction(ctx, func(ctx context.Context) error {
		var err error
		user, err := s.Repository.GetByID(ctx, id)
		if err != nil {
			return errors.New(err)
		}

		user = typeutil.Copy(user, updateDTO)
		if updateDTO.Avatar.Present {
			// Delete previous avatar
			if user.Avatar.Valid {
				err := s.StorageService.Delete(user.Avatar.String)
				if err != nil {
					return errors.New(err)
				}
			}
			if updateDTO.Avatar.Val.Valid {
				// Save new avatar
				filename, err := s.StorageService.SaveAvatar(updateDTO.Avatar.Val.V[0])
				if err != nil {
					return errors.New(err)
				}
				user.Avatar.SetValid(filename)
			} else {
				user.Avatar = null.String{}
			}
		}

		_, err = s.Repository.Update(ctx, user)
		return errors.New(err)
	})

	return errors.New(err)
}

func (s *Service) Name() string {
	return service.User
}
