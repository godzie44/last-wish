package application

import (
	"context"
	"errors"
	"lw/internal/domain"
)

type UserService struct {
	repo     domain.UserRepository
	notifier domain.NotifyService
}

func NewUserService(repo domain.UserRepository, notifier domain.NotifyService) *UserService {
	return &UserService{repo: repo, notifier: notifier}
}

var ErrEmailNotUnique = errors.New("email must be")

func (u *UserService) NewUser(ctx context.Context, name, email string) error {
	if _, err := u.repo.FindByEmail(ctx, email); err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			return err
		}
	} else {
		return ErrEmailNotUnique
	}

	user, err := domain.NewUser(name, email)
	if err != nil {
		return err
	}

	return u.repo.Persists(ctx, user)
}

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

func (u *UserService) ReleaseWishes(ctx context.Context, userId int64) error {
	user, err := u.repo.FindById(ctx, userId)
	if err != nil {
		return err
	}

	user.ReleaseWishes(u.notifier)
	return nil
}
