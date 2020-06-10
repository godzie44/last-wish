package application

import (
	"context"
	"github.com/stretchr/testify/assert"
	"lw/internal/domain"
	"testing"
)

func TestNewWish(t *testing.T) {
	repo := &inMemoryUserRepo{}
	user, _ := domain.NewUser("test", "test@test.com")
	assert.NoError(t, repo.Persists(context.Background(), user))
	service := NewWishService(repo)
	assert.NoError(t, service.NewWish(context.Background(), 0, "some title"))
}

func TestNewWishFailIfUserNotFound(t *testing.T) {
	repo := &inMemoryUserRepo{}
	service := NewWishService(repo)
	assert.Error(t, service.NewWish(context.Background(), 0, "some title"))
}

type inMemoryUserRepo struct {
	users []*domain.User
}

func (i *inMemoryUserRepo) Persists(_ context.Context, u *domain.User) error {
	i.users = append(i.users, u)
	return nil
}

func (i *inMemoryUserRepo) FindById(_ context.Context, id int64) (*domain.User, error) {
	if int(id) >= len(i.users) {
		return nil, domain.ErrUserNotFound
	}
	return i.users[int(id)], nil
}

func (i *inMemoryUserRepo) FindByEmail(_ context.Context, mail string) (*domain.User, error) {
	for _, user := range i.users {
		if user.Email() == mail {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}
