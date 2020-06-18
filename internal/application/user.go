package application

import (
	"context"
	"errors"
	"github.com/godzie44/lw/internal/domain"
)

// UserService holds methods for work with user aggregate: add friends, create users, release wishes.
type UserService struct {
	repo     domain.UserRepository
	notifier domain.NotifyService
}

func NewUserService(repo domain.UserRepository, notifier domain.NotifyService) *UserService {
	return &UserService{repo: repo, notifier: notifier}
}

var ErrEmailNotUnique = errors.New("email must be unique")

// NewUser create new user.
func (u *UserService) NewUser(ctx context.Context, name, email string) (int64, error) {
	if _, err := u.repo.FindByEmail(ctx, email); err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			return 0, err
		}
	} else {
		return 0, ErrEmailNotUnique
	}

	user, err := domain.NewUser(name, email)
	if err != nil {
		return 0, err
	}

	if err = u.repo.Add(ctx, user); err != nil {
		return 0, err
	}

	return user.ID(), err
}

// AddFriend add friend relation between two users.
func (u *UserService) AddFriend(ctx context.Context, userId, friendId int64) error {
	user, err := u.repo.FindById(ctx, userId)
	if err != nil {
		return err
	}

	friend, err := u.repo.FindById(ctx, friendId)
	if err != nil {
		return err
	}

	return user.AddFriend(friend)
}

// ReleaseWishes release all user wishes.
func (u *UserService) ReleaseWishes(ctx context.Context, userId int64) error {
	user, err := u.repo.FindById(ctx, userId)
	if err != nil {
		return err
	}

	user.ReleaseWishes(u.notifier)
	return nil
}
