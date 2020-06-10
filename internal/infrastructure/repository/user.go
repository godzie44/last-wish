package repository

import (
	"context"
	"github.com/godzie44/d3/orm"
	"lw/internal/domain"
)

type D3UserRepo struct {
	d3Rep *orm.Repository
}

func NewD3UserRepo(d3Rep *orm.Repository) *D3UserRepo {
	return &D3UserRepo{d3Rep: d3Rep}
}

func (d *D3UserRepo) Persists(ctx context.Context, u *domain.User) error {
	return d.d3Rep.Persists(ctx, u)
}

func (d *D3UserRepo) FindById(ctx context.Context, id int64) (*domain.User, error) {
	e, err := d.d3Rep.FindOne(ctx, d.d3Rep.CreateQuery().AndWhere("id=?", id))
	if err != nil {
		return nil, err
	}

	return e.(*domain.User), nil
}

func (d *D3UserRepo) FindByEmail(ctx context.Context, mail string) (*domain.User, error) {
	e, err := d.d3Rep.FindOne(ctx, d.d3Rep.CreateQuery().AndWhere("email=?", mail))
	if err != nil {
		if err == orm.ErrEntityNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return e.(*domain.User), nil
}
