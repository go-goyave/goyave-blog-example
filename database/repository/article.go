package repository

import (
	"context"

	"github.com/go-goyave/goyave-blog-example/database/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"goyave.dev/filter"
	"goyave.dev/goyave/v5/database"
	"goyave.dev/goyave/v5/util/errors"
	"goyave.dev/goyave/v5/util/session"
	"goyave.dev/goyave/v5/util/typeutil"
)

type Article struct {
	DB *gorm.DB
}

func NewArticle(db *gorm.DB) *Article {
	return &Article{
		DB: db,
	}
}

func (r *Article) Index(ctx context.Context, request *filter.Request) (*database.Paginator[*model.Article], error) {
	settings := &filter.Settings[*model.Article]{
		DefaultSort: []*filter.Sort{
			{Field: "created_at", Order: filter.SortDescending},
		},
		FieldsSearch: []string{"title"},
		Blacklist: filter.Blacklist{
			FieldsBlacklist: []string{"deleted_at"},
			Relations: map[string]*filter.Blacklist{
				"Author": {IsFinal: true},
			},
		},
	}
	paginator, err := settings.Scope(session.DB(ctx, r.DB), request, &[]*model.Article{})
	return paginator, errors.New(err)
}

func (r *Article) GetByID(ctx context.Context, id int64) (*model.Article, error) {
	var article *model.Article
	db := session.DB(ctx, r.DB).Where("id", id).First(&article)
	return article, errors.New(db.Error)
}

func (r *Article) GetBySlug(ctx context.Context, slug string) (*model.Article, error) {
	var article *model.Article
	db := session.DB(ctx, r.DB).Where("slug", slug).First(&article)
	return article, errors.New(db.Error)
}

func (r *Article) Create(ctx context.Context, article *model.Article) (*model.Article, error) {
	db := session.DB(ctx, r.DB).Omit(clause.Associations).Create(&article)
	return article, errors.New(db.Error)
}

func (r *Article) Update(ctx context.Context, article *model.Article) (*model.Article, error) {
	if article.ID.Val == 0 {
		return nil, errors.New(gorm.ErrPrimaryKeyRequired)
	}
	db := session.DB(ctx, r.DB).Omit(clause.Associations).Save(&article)
	return article, errors.New(db.Error)
}

func (r *Article) Delete(ctx context.Context, id int64) error {
	db := session.DB(ctx, r.DB).Delete(&model.Article{ID: typeutil.NewUndefined(id)})
	if db.RowsAffected == 0 {
		return errors.New(gorm.ErrRecordNotFound)
	}
	return errors.New(db.Error)
}

func (r *Article) IsOwner(ctx context.Context, resourceID, ownerID int64) (bool, error) {
	var one int64
	db := session.DB(ctx, r.DB).
		Table(model.Article{}.TableName()).
		Select("1").
		Where("id", resourceID).
		Where("author_id", ownerID).
		Where("deleted_at IS NULL").
		Find(&one)
	return one == 1, errors.New(db.Error)
}
