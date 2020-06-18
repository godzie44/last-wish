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
