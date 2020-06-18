package application

import (
	"context"
	"github.com/godzie44/lw/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUser(t *testing.T) {
	repo := &inMemoryUserRepo{}
	service := NewUserService(repo, nil)

	_, err := service.NewUser(context.Background(), "test", "test@test.com")
	assert.NoError(t, err)
	assert.Len(t, repo.users, 1)
}

func TestNewUserFailIfEmailNotUnique(t *testing.T) {
	repo := &inMemoryUserRepo{}
	u, _ := domain.NewUser("", "test@test.com")
	assert.NoError(t, repo.Add(context.Background(), u))
	service := NewUserService(repo, nil)

	_, err := service.NewUser(context.Background(), "test", "test@test.com")
	assert.EqualError(t, err, ErrEmailNotUnique.Error())
}

func TestAddFriend(t *testing.T) {
	repo := &inMemoryUserRepo{}
	u1, _ := domain.NewUser("u1", "test@test.com")
	u2, _ := domain.NewUser("u2", "test2@test.com")
	assert.NoError(t, repo.Add(context.Background(), u1))
	assert.NoError(t, repo.Add(context.Background(), u2))
	service := NewUserService(repo, nil)

	assert.NoError(t, service.AddFriend(context.Background(), 0, 1))

	assert.Equal(t, 1, u1.FriendsCount())
}

func TestAddFriendFailIfUserNotExists(t *testing.T) {
	repo := &inMemoryUserRepo{}
	u1, _ := domain.NewUser("u1", "test@test.com")
	assert.NoError(t, repo.Add(context.Background(), u1))
	service := NewUserService(repo, nil)

	assert.EqualError(t, service.AddFriend(context.Background(), 0, 1), domain.ErrUserNotFound.Error())
}

func TestReleaseWishesFailIfUserNotExists(t *testing.T) {
	repo := &inMemoryUserRepo{}
	service := NewUserService(repo, nil)

	assert.EqualError(t, service.ReleaseWishes(context.Background(), 0), domain.ErrUserNotFound.Error())
}
