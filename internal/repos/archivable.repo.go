package repos

import (
	"context"

	"gorm.io/gorm"
)

type archivable[T any] struct {
	db *gorm.DB
}

func (r *archivable[T]) Archive(ctx context.Context, where any, args ...any) (err error) {
	result, err := gorm.G[T](r.db).Where(where, args...).Delete(ctx)
	if result == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}
	return
}

func (r *archivable[T]) Restore(ctx context.Context, where any, args ...any) (err error) {
	result, err := gorm.G[T](r.db.Unscoped()).Where(where, args...).Update(ctx, "deleted_at", nil)
	if result == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}
	return
}
