package repository

import (
	"context"
	"github.com/godzie44/d3/orm"
	"github.com/godzie44/lw/internal/domain"
)

type D3UserRepo struct {
	d3Repo *orm.Repository
}

func NewUserRepository(d3Rep *orm.Repository) *D3UserRepo {
	return &D3UserRepo{d3Repo: d3Rep}
}

func (d *D3UserRepo) Add(ctx context.Context, u *domain.User) error {
	if err := d.d3Repo.Persists(ctx, u); err != nil {
		return err
	}

	return orm.Session(ctx).Flush()
}

func (d *D3UserRepo) FindById(ctx context.Context, id int64) (*domain.User, error) {
	e, err := d.d3Repo.FindOne(ctx, d.d3Repo.Select().AndWhere("id", "=", id))
	if err != nil {
		return nil, err
	}

	return e.(*domain.User), nil
}

func (d *D3UserRepo) FindByEmail(ctx context.Context, mail string) (*domain.User, error) {
	e, err := d.d3Repo.FindOne(ctx, d.d3Repo.Select().AndWhere("email", "=", mail))
	if err != nil {
		if err == orm.ErrEntityNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return e.(*domain.User), nil
}
