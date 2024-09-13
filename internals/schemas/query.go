package schemas

import (
	"context"
	"errors"
	"fmt"

	"github.com/thrillee/triq/internals/common"
	"gorm.io/gorm"
)

type Model interface {
	// ToString() string
	GetID() interface{}
	GetModelName() interface{}
}

func convertModel[T any](model Model) (T, bool) {
	todo, ok := model.(T)
	return todo, ok
}

type QueryResult struct {
	Results *gorm.DB
	Page    common.Pagination
}

type Repository[T Model] interface {
	Exists(context.Context, *gorm.DB) bool
	GetByID(context.Context, interface{}) (*T, error)
	Filter(context.Context, Query) *QueryResult
	Create(context.Context, *T) (*T, error)
	Edit(context.Context, interface{}, *T) (*T, error)
	Delete(context.Context, *T) error
	Query(context.Context) *gorm.DB
	Save(context.Context, *T) (*T, error)
	Search(context.Context, []string, string, common.Paginable) *gorm.DB
	GetModel() *Model
}

type baseRepository[T Model] struct {
	db    *gorm.DB
	model Model
}

func (b *baseRepository[T]) GetModel() *Model {
	return &b.model
}

type Query interface {
	GetLimit() int
	GetOffset() int
	GetQuerySet() *gorm.DB
	GetCurrentURL() string
}

type DefaultQuery struct {
	CurrentURL string
	Limit      int
	Offset     int
	QuerySet   *gorm.DB
}

func (dq DefaultQuery) GetLimit() int {
	return dq.Limit
}

func (dq DefaultQuery) GetOffset() int {
	return dq.Offset
}

func (dq DefaultQuery) GetQuerySet() *gorm.DB {
	return dq.QuerySet
}

func (dq DefaultQuery) GetCurrentURL() string {
	return dq.CurrentURL
}

func NewBaseRepository[T Model](model Model) Repository[T] {
	return &baseRepository[T]{db: getDB(), model: model}
}

func (repo baseRepository[T]) Exists(ctx context.Context, query *gorm.DB) bool {
	var count int64
	query.WithContext(ctx).Count(&count)
	return count > 0
}

func (repo baseRepository[T]) Query(ctx context.Context) *gorm.DB {
	return repo.db.WithContext(ctx).Model(repo.model)
}

func (repo baseRepository[T]) GetByID(ctx context.Context, modelId interface{}) (*T, error) {
	var model T
	result := db.WithContext(ctx).Where("id = ?", modelId).First(&model)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &model, nil
}

func (repo baseRepository[T]) Search(ctx context.Context, search []string, searchValue string, page common.Paginable) *gorm.DB {
	query := repo.Query(ctx)

	searchCount := len(search)
	count := 0
	for _, k := range search {
		likeValue := fmt.Sprintf("%%%s%%", searchValue)

		if count == 0 {
			query = query.Where("?=?", k, likeValue)
		}

		if count > searchCount {
			query = query.Or("?=?", k, likeValue)
		}

		count += 1
	}

	return query
}

func (repo baseRepository[T]) Filter(ctx context.Context, queriable Query) *QueryResult {
	query := queriable.GetQuerySet().WithContext(ctx)

	var totalRecords int64
	query.Count(&totalRecords)

	if queriable.GetLimit() > 0 && queriable.GetOffset() >= 0 {
		query = query.Limit(queriable.GetLimit()).Offset(queriable.GetOffset())
	}

	return &QueryResult{
		Results: query,
		Page: common.CreatePagination(&common.PageParams{
			Count:  totalRecords,
			Limit:  queriable.GetLimit(),
			Offset: queriable.GetOffset(),
			URL:    queriable.GetCurrentURL(),
		}),
	}
}

func (repo baseRepository[T]) Delete(ctx context.Context, model *T) error {
	result := repo.db.WithContext(ctx).Delete(model)
	return result.Error
}

func (repo baseRepository[T]) Edit(ctx context.Context, id interface{}, model *T) (*T, error) {
	result := repo.Query(ctx).Where("id = ?", id).Updates(model)
	return model, result.Error
}

func (repo baseRepository[T]) Create(ctx context.Context, model *T) (*T, error) {
	result := repo.db.WithContext(ctx).Create(model)
	return model, result.Error
}

func (repo baseRepository[T]) Save(ctx context.Context, model *T) (*T, error) {
	repo.Query(ctx).Save(model)
	return model, nil
}
