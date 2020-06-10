package application

import (
	"context"
	"lw/internal/domain"
)

type WishService struct {
	userRepo domain.UserRepository
}

func (w *WishService) NewWish(ctx context.Context, userId int64, content string) error {
	user, err := w.userRepo.FindById(ctx, userId)
	if err != nil {
		return err
	}

	return user.AddWish(content)
}
