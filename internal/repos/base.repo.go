package repos

import (
	"context"
	"school-information-system/config"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type create[T any] struct {
	db *gorm.DB
}

type read[T any] struct {
	db *gorm.DB
}

type update[T any] struct {
	db *gorm.DB
}

type delete[T any] struct {
	db *gorm.DB
}

// create

func (r *create[T]) Create(ctx context.Context, data *T) error {
	return gorm.G[T](r.db).Create(ctx, data)
}

func (r *create[T]) CreateAll(ctx context.Context, data *[]T) error {
	return gorm.G[T](r.db).CreateInBatches(ctx, data, config.MAX_CREATE_BATCH)
}

// read

func (r *read[T]) GetFirst(ctx context.Context, where any, args ...any) (T, error) {
	return gorm.G[T](r.db).Where(where, args...).First(ctx)
}
func (r *read[T]) GetByID(ctx context.Context, id string) (T, error) {
	return r.GetFirst(ctx, "id = ?", id)
}

func (r *read[T]) GetAll(ctx context.Context, where any, args ...any) ([]T, error) {
	return gorm.G[T](r.db).Where(where, args...).Find(ctx)
}

func (r *read[T]) GetByIDs(ctx context.Context, ids []string) ([]T, error) {
	return r.GetAll(ctx, "id IN ?", ids)
}

// update

func (r *update[T]) Update(ctx context.Context, u T, where any, args ...any) error {
	_, err := gorm.G[T](r.db).Where(where, args...).Updates(ctx, u)
	return err
}

func (r *update[T]) UpdateByID(ctx context.Context, id string, u T) error {
	return r.Update(ctx, u, "id = ?", id)
}

func (r *update[T]) UpdateAndGet(ctx context.Context, u T, where any, args ...any) (result T, err error) {
	err = r.db.WithContext(ctx).Model(new(T)).Clauses(clause.Returning{}).Where(where, args...).Updates(u).Scan(&result).Error
	return
}

func (r *update[T]) UpdateByIDAndGet(ctx context.Context, id string, u T) (result T, err error) {
	return r.UpdateAndGet(ctx, u, "id = ?", id)
}

// delete

func (r *delete[T]) Delete(ctx context.Context, where any, args ...any) error {
	_, err := gorm.G[T](r.db).Where(where, args...).Delete(ctx)
	return err
}

func (r *delete[T]) DeleteByID(ctx context.Context, id string) error {
	return r.Delete(ctx, "id = ?", id)
}

func (r *delete[T]) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.Delete(ctx, "id IN ?", ids)
}
