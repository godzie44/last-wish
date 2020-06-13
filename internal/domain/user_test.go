package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeWish(t *testing.T) {
	user, _ := NewUser("test", "test@test.com")

	assert.NoError(t, user.MakeWish("test_wish_1"))
}

func TestAddWishAfterLimitExceeded(t *testing.T) {
	user, _ := NewUser("test", "test@test.com")

	for i := 0; i < maxWishesPerUser; i++ {
		assert.NoError(t, user.MakeWish("test_wish"))
	}

	assert.Equal(t, maxWishesPerUser, user.wishes.Count())
	assert.EqualError(t, user.MakeWish("test_wish"), errToManyWishes.Error())
	assert.Equal(t, maxWishesPerUser, user.wishes.Count())
}

func TestAddFriend(t *testing.T) {
	user, _ := NewUser("test", "test@test.com")
	friend, _ := NewUser("friend", "friend@test.com")

	assert.NoError(t, user.AddFriend(friend))
}

func TestAddFriendAfterLimitExceeded(t *testing.T) {
	user, _ := NewUser("test", "test@test.com")

	for i := 0; i < maxFriendsPerUser; i++ {
		friend, _ := NewUser("friend", "friend@test.com")
		assert.NoError(t, user.AddFriend(friend))
	}

	assert.Equal(t, maxFriendsPerUser, user.friends.Count())
	friend, _ := NewUser("friend", "friend@test.com")
	assert.EqualError(t, user.AddFriend(friend), ErrToManyFriends.Error())
	assert.Equal(t, maxFriendsPerUser, user.friends.Count())
}

func TestCantCreateUserWithLongName(t *testing.T) {
	_, err := NewUser("veryveryveryveryveryveryveryveryvery long name", "some@mail.com")
	assert.EqualError(t, err, errUserNameTooLong.Error())
}
