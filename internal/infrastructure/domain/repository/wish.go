package repository

import (
	"context"
	"github.com/godzie44/d3/orm"
	"github.com/godzie44/lw/internal/domain"
)

type d3WishRepository struct {
	d3Rep *orm.Repository
}

func NewWishRepository(d3Rep *orm.Repository) *d3WishRepository {
	return &d3WishRepository{d3Rep: d3Rep}
}

func (d *d3WishRepository) Delete(ctx context.Context, wish *domain.Wish) error {
	return d.d3Rep.Delete(ctx, wish)
}

func (d *d3WishRepository) FindByUserIdAndId(ctx context.Context, userId, wishId int64) (*domain.Wish, error) {
	w, err := d.d3Rep.FindOne(ctx, d.d3Rep.Select().Where("user_id", "=", userId).AndWhere("id", "=", wishId))
	if err != nil {
		return nil, err
	}
	return w.(*domain.Wish), nil
}
