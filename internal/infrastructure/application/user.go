package application

import (
	"context"
	"github.com/godzie44/d3/orm"
	"github.com/godzie44/lw/internal/application"
)

type transactionalUserService struct {
	next *application.UserService
}

func NewTransactionalUserService(us *application.UserService) *transactionalUserService {
	return &transactionalUserService{next: us}
}

func (t *transactionalUserService) NewUser(ctx context.Context, name, email string) (int64, error) {
	return t.next.NewUser(ctx, name, email)
}

func (t *transactionalUserService) AddFriend(ctx context.Context, userId, friendId int64) error {
	if err := t.next.AddFriend(ctx, userId, friendId); err != nil {
		return err
	}
	return orm.Session(ctx).Flush()
}

// ReleaseWishes release all user wishes.
func (t *transactionalUserService) ReleaseWishes(ctx context.Context, userId int64) error {
	if err := t.next.ReleaseWishes(ctx, userId); err != nil {
		return err
	}
	return orm.Session(ctx).Flush()
}
