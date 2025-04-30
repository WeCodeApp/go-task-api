package core

import (
	"context"
	"gorm.io/gorm"
)

var (
	ErrNotFound = gorm.ErrRecordNotFound
)

type Logger interface {
	Error(args ...interface{})
}

type Service interface {
	HandleError(ctx context.Context, err error) error
}

type GormRepository interface {
	GetAll(ctx context.Context, target interface{}, preloads ...string) error
	GetAllOrdered(ctx context.Context, target interface{}, order string, preloads ...string) error
	GetAllCount(ctx context.Context, target interface{}, count *int64, preloads ...string) error

	GetBatch(ctx context.Context, target interface{}, limit, offset int, preloads ...string) error
	GetBatchOrdered(ctx context.Context, target interface{}, limit, offset int, order string, preloads ...string) error
	GetBatchCount(ctx context.Context, target interface{}, count *int64, limit, offset int, preloads ...string) error

	GetWhere(ctx context.Context, target interface{}, condition string, preloads ...string) error
	GetWhereOrdered(ctx context.Context, target interface{}, condition, order string, preloads ...string) error
	GetWhereCount(ctx context.Context, target interface{}, count *int64, condition string, preloads ...string) error

	GetWhereBatch(ctx context.Context, target interface{}, condition string, limit, offset int, preloads ...string) error
	GetWhereBatchOrdered(ctx context.Context, target interface{}, condition string, limit, offset int, order string, preloads ...string) error
	GetWhereBatchCount(ctx context.Context, target interface{}, count *int64, condition string, limit, offset int, preloads ...string) error

	GetByField(ctx context.Context, target interface{}, field string, value interface{}, preloads ...string) error
	GetByFieldWithInterval(ctx context.Context, target interface{}, field string, value interface{}, preloads ...string) error
	GetByFieldOrdered(ctx context.Context, target interface{}, field string, value interface{}, order string, preloads ...string) error
	GetByFieldCount(ctx context.Context, target interface{}, count *int64, field string, value interface{}, preloads ...string) error

	GetByFields(ctx context.Context, target interface{}, filters map[string]interface{}, preloads ...string) error
	GetByFieldsOrdered(ctx context.Context, target interface{}, filters map[string]interface{}, order string, preloads ...string) error
	GetByFieldsCount(ctx context.Context, target interface{}, count *int64, filters map[string]interface{}, preloads ...string) error

	GetByFieldBatch(ctx context.Context, target interface{}, field string, value interface{}, limit, offset int, preloads ...string) error
	GetByFieldBatchOrdered(ctx context.Context, target interface{}, field string, value interface{}, limit, offset int, order string, preloads ...string) error
	GetByFieldBatchCount(ctx context.Context, target interface{}, count *int64, field string, value interface{}, limit, offset int, preloads ...string) error

	GetByFieldsBatch(ctx context.Context, target interface{}, filters map[string]interface{}, limit, offset int, preloads ...string) error
	GetByFieldsBatchOrdered(ctx context.Context, target interface{}, filters map[string]interface{}, limit, offset int, order string, preloads ...string) error
	GetByFieldsBatchCount(ctx context.Context, target interface{}, count *int64, filters map[string]interface{}, limit, offset int, preloads ...string) error

	GetOneByField(ctx context.Context, target interface{}, field string, value interface{}, preloads ...string) error
	GetOneByFieldOrdered(ctx context.Context, target interface{}, field string, value interface{}, order string, preloads ...string) error
	GetOneByFields(ctx context.Context, target interface{}, filters map[string]interface{}, preloads ...string) error
	GetOneByFieldsOrdered(ctx context.Context, target interface{}, filters map[string]interface{}, order string, preloads ...string) error

	// GetOneByID assumes you have a PK column "id". If this is not the case just ignore the method
	// and add a custom struct with this GormRepository embedded.
	GetOneByID(ctx context.Context, target interface{}, id string, preloads ...string) error

	Create(ctx context.Context, target interface{}) error
	Save(ctx context.Context, target interface{}) error
	Delete(ctx context.Context, target interface{}) error
	DeleteByField(ctx context.Context, target interface{}, field string, value interface{}) error

	DB() *gorm.DB
	HandleError(context.Context, *gorm.DB) error
	HandleOneError(context.Context, *gorm.DB) error
}

// GormTransactionRepository extends GormRepository with modifier functions that accept a transaction
type GormTransactionRepository interface {
	GormRepository
	CreateTx(ctx context.Context, target interface{}, tx *gorm.DB) error
	SaveTx(ctx context.Context, target interface{}, tx *gorm.DB) error
	DeleteTx(ctx context.Context, target interface{}, tx *gorm.DB) error
}
