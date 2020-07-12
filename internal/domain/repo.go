package domain

import (
	"context"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Add(ctx context.Context, u *User) error
	FindById(ctx context.Context, id int64) (*User, error)
	FindByEmail(ctx context.Context, mail string) (*User, error)
}

type WishRepository interface {
	Delete(ctx context.Context, wish *Wish) error
	FindByUserIdAndId(ctx context.Context, userId, wishId int64) (*Wish, error)
}
