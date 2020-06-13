package application

import (
	"context"
	"lw/internal/domain"
)

// WishService holds methods for work with wishes.
type WishService struct {
	userRepo domain.UserRepository
}

func NewWishService(userRepo domain.UserRepository) *WishService {
	return &WishService{userRepo: userRepo}
}

// NewWish create new wish.
func (w *WishService) NewWish(ctx context.Context, userId int64, content string) error {
	user, err := w.userRepo.FindById(ctx, userId)
	if err != nil {
		return err
	}

	return user.MakeWish(content)
}
