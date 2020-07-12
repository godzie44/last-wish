package application

import (
	"context"
	"github.com/godzie44/lw/internal/domain"
)

// WishService holds methods for work with wishes.
type WishService struct {
	userRepo domain.UserRepository
	wishRepo domain.WishRepository
}

func NewWishService(userRepo domain.UserRepository, wishRepo domain.WishRepository) *WishService {
	return &WishService{userRepo: userRepo, wishRepo: wishRepo}
}

// NewWish create new wish.
func (w *WishService) NewWish(ctx context.Context, userId int64, content string) error {
	user, err := w.userRepo.FindById(ctx, userId)
	if err != nil {
		return err
	}

	return user.MakeWish(content)
}

func (w *WishService) UpdateWish(ctx context.Context, userId int64, wishId int64, content string) error {
	user, err := w.userRepo.FindById(ctx, userId)
	if err != nil {
		return err
	}

	return user.UpdateWish(wishId, content)
}

func (w *WishService) DeleteWish(ctx context.Context, userId int64, wishId int64) error {
	wish, err := w.wishRepo.FindByUserIdAndId(ctx, userId, wishId)
	if err != nil {
		return err
	}

	return w.wishRepo.Delete(ctx, wish)
}
