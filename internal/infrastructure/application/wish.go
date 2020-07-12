package application

import (
	"context"
	"github.com/godzie44/d3/orm"
	"github.com/godzie44/lw/internal/application"
)

type transactionalWishService struct {
	next *application.WishService
}

func NewTransactionalWishService(next *application.WishService) *transactionalWishService {
	return &transactionalWishService{next: next}
}

func (t *transactionalWishService) NewWish(ctx context.Context, userId int64, content string) error {
	if err := t.next.NewWish(ctx, userId, content); err != nil {
		return err
	}
	return orm.Session(ctx).Flush()
}

func (t *transactionalWishService) UpdateWish(ctx context.Context, userId int64, wishId int64, content string) error {
	if err := t.next.UpdateWish(ctx, userId, wishId, content); err != nil {
		return err
	}
	return orm.Session(ctx).Flush()
}

func (t *transactionalWishService) DeleteWish(ctx context.Context, userId int64, wishId int64) error {
	if err := t.next.DeleteWish(ctx, userId, wishId); err != nil {
		return err
	}
	return orm.Session(ctx).Flush()
}
