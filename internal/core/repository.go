package core

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type repository struct {
	logger       Logger
	db           *gorm.DB
	defaultJoins []string
}

func NewRepository(db *gorm.DB, defaultJoins ...string) GormTransactionRepository {
	return &repository{
		defaultJoins: defaultJoins,
		db:           db,
	}
}

func (r *repository) DB() *gorm.DB {
	return r.WithPreloads(nil)
}

func (r *repository) GetAll(ctx context.Context, target interface{}, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Unscoped().
		Find(target)

	return r.HandleError(ctx, res)
}

func (r *repository) GetAllOrdered(ctx context.Context, target interface{}, order string, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Unscoped().
		Order(order).
		Find(target)

	return r.HandleError(ctx, res)
}

func (r *repository) GetAllCount(ctx context.Context, target interface{}, count *int64, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Unscoped().
		Find(target).
		Count(count)

	return r.HandleError(ctx, res)
}

func (r *repository) GetBatch(ctx context.Context, target interface{}, limit, offset int, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Unscoped().
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(ctx, res)
}

func (r *repository) GetBatchOrdered(ctx context.Context, target interface{}, limit, offset int, order string, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Unscoped().
		Limit(limit).
		Offset(offset).
		Order(order).
		Find(target)

	return r.HandleError(ctx, res)
}

func (r *repository) GetBatchCount(ctx context.Context, target interface{}, count *int64, limit, offset int, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Unscoped().
		Limit(limit).
		Offset(offset).
		Find(target).
		Count(count)

	return r.HandleError(ctx, res)
}

func (r *repository) GetWhere(ctx context.Context, target interface{}, condition string, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Where(condition).
		Find(target)

	return r.HandleError(ctx, res)
}

func (r *repository) GetWhereOrdered(ctx context.Context, target interface{}, condition, order string, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Where(condition).
		Order(order).
		Find(target)

	return r.HandleError(ctx, res)
}

func (r *repository) GetWhereCount(ctx context.Context, target interface{}, count *int64, condition string, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Where(condition).
		Find(target).
		Count(count)

	return r.HandleError(ctx, res)
}

func (r *repository) GetWhereBatch(ctx context.Context, target interface{}, condition string, limit, offset int, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Where(condition).
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(ctx, res)
}

func (r *repository) GetWhereBatchOrdered(ctx context.Context, target interface{}, condition string, limit, offset int, order string, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Where(condition).
		Limit(limit).
		Offset(offset).
		Order(order).
		Find(target)

	return r.HandleError(ctx, res)
}

func (r *repository) GetWhereBatchCount(ctx context.Context, target interface{}, count *int64, condition string, limit, offset int, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Where(condition).
		Limit(limit).
		Offset(offset).
		Find(target).
		Count(count)

	return r.HandleError(ctx, res)
}

func (r *repository) GetByField(ctx context.Context, target interface{}, field string, value interface{}, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		Find(target)

	return r.HandleError(ctx, res)
}

func (r *repository) GetByFieldOrdered(ctx context.Context, target interface{}, field string, value interface{}, order string, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		Order(order).
		Find(target)

	return r.HandleError(ctx, res)
}

func (r *repository) GetByFieldCount(ctx context.Context, target interface{}, count *int64, field string, value interface{}, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		Find(target).
		Count(count)

	return r.HandleError(ctx, res)
}

func (r *repository) GetByFields(ctx context.Context, target interface{}, filters map[string]interface{}, preloads ...string) error {
	db := r.WithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	return r.HandleError(ctx, db.Find(target))
}

func (r *repository) GetByFieldWithInterval(ctx context.Context, target interface{}, field string, value interface{}, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		Where("time_out > now() - interval 1 hour").
		Find(target)

	return r.HandleOneError(ctx, res.First(target))
}

func (r *repository) GetByFieldsOrdered(ctx context.Context, target interface{}, filters map[string]interface{}, order string, preloads ...string) error {
	db := r.WithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	return r.HandleError(ctx, db.Order(order).Find(target))
}

func (r *repository) GetByFieldsCount(ctx context.Context, target interface{}, count *int64, filters map[string]interface{}, preloads ...string) error {
	db := r.WithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v ?", field), value)
	}

	return r.HandleError(ctx, db.Find(target).Count(count))
}

func (r *repository) GetByFieldBatch(ctx context.Context, target interface{}, field string, value interface{}, limit, offset int, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(ctx, res)
}

func (r *repository) GetByFieldBatchOrdered(ctx context.Context, target interface{}, field string, value interface{}, limit, offset int, order string, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		Limit(limit).
		Offset(offset).
		Order(order).
		Find(target)

	return r.HandleError(ctx, res)
}

func (r *repository) GetByFieldBatchCount(ctx context.Context, target interface{}, count *int64, field string, value interface{}, limit, offset int, preloads ...string) error {
	res := r.WithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		Limit(limit).
		Offset(offset).
		Find(target).
		Count(count)

	return r.HandleError(ctx, res)
}

func (r *repository) GetByFieldsBatch(ctx context.Context, target interface{}, filters map[string]interface{}, limit, offset int, preloads ...string) error {
	db := r.WithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	res := db.
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(ctx, res)
}

func (r *repository) GetByFieldsBatchOrdered(ctx context.Context, target interface{}, filters map[string]interface{}, limit, offset int, order string, preloads ...string) error {
	db := r.WithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	res := db.
		Limit(limit).
		Offset(offset).
		Order(order).
		Find(target)

	return r.HandleError(ctx, res)
}

func (r *repository) GetByFieldsBatchCount(ctx context.Context, target interface{}, count *int64, filters map[string]interface{}, limit, offset int, preloads ...string) error {
	db := r.WithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	res := db.
		Limit(limit).
		Offset(offset).
		Find(target).
		Count(count)

	return r.HandleError(ctx, res)
}

func (r *repository) GetOneByField(ctx context.Context, target interface{}, field string, value interface{}, preloads ...string) error {
	return r.HandleOneError(ctx, r.WithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		First(target))
}

func (r *repository) GetOneByFieldOrdered(ctx context.Context, target interface{}, field string, value interface{}, order string, preloads ...string) error {
	return r.HandleOneError(ctx, r.WithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		Order(order).
		First(target))
}

func (r *repository) GetOneByFields(ctx context.Context, target interface{}, filters map[string]interface{}, preloads ...string) error {
	db := r.WithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v ?", field), value)
	}

	return r.HandleOneError(ctx, db.First(target))
}

func (r *repository) GetOneByFieldsOrdered(ctx context.Context, target interface{}, filters map[string]interface{}, order string, preloads ...string) error {
	db := r.WithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	return r.HandleOneError(ctx, db.Order(order).First(target))
}

func (r *repository) GetOneByID(ctx context.Context, target interface{}, id string, preloads ...string) error {
	return r.HandleOneError(ctx, r.WithPreloads(preloads).
		Where("id = ?", id).
		First(target))
}

func (r *repository) Create(ctx context.Context, target interface{}) error {
	return r.HandleError(ctx, r.db.Create(target))
}

func (r *repository) CreateTx(ctx context.Context, target interface{}, tx *gorm.DB) error {
	return r.HandleError(ctx, tx.Create(target))
}

func (r *repository) Save(ctx context.Context, target interface{}) error {
	return r.HandleError(ctx, r.db.Save(target))
}

func (r *repository) SaveTx(ctx context.Context, target interface{}, tx *gorm.DB) error {
	return r.HandleError(ctx, tx.Save(target))
}

func (r *repository) Delete(ctx context.Context, target interface{}) error {
	return r.HandleError(ctx, r.db.Delete(target))
}

func (r *repository) DeleteByField(ctx context.Context, target interface{}, field string, value interface{}) error {
	res := r.db.Where(fmt.Sprintf("%v = ?", field), value).
		Delete(target)

	return r.HandleError(ctx, res)
}

func (r *repository) DeleteTx(ctx context.Context, target interface{}, tx *gorm.DB) error {
	return r.HandleError(ctx, tx.Delete(target))
}

func (r *repository) HandleError(ctx context.Context, res *gorm.DB) error {
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		r.logger.Error(res.Error)
		return res.Error
	}

	return nil
}

func (r *repository) HandleOneError(ctx context.Context, res *gorm.DB) error {
	if err := r.HandleError(ctx, res); err != nil {
		return err
	}

	if res.RowsAffected != 1 {
		//r.logger.Error(ErrNotFound)
		return ErrNotFound
	}

	return nil
}

func (r *repository) WithPreloads(preloads []string) *gorm.DB {
	dbConn := r.db

	for _, join := range r.defaultJoins {
		dbConn = dbConn.Joins(join)
	}

	for _, preload := range preloads {
		dbConn = dbConn.Preload(preload)
	}

	return dbConn
}
