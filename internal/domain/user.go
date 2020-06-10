package domain

import (
	"database/sql"
	"errors"
	"github.com/godzie44/d3/orm/entity"
)

//d3:entity
//d3_table:lw_user
type User struct {
	id      sql.NullInt64      `d3:"pk:auto"`
	name    string             `d3:"column:name"`
	email   string             `d3:"column:email"`
	wishes  *entity.Collection `d3:"one_to_many:<target_entity:lw/internal/domain/Wish,join_on:user_id,delete:nullable>"`
	friends *entity.Collection `d3:"many_to_many:<target_entity:lw/internal/domain/User,join_on:u1_id,reference_on:u2_id,join_table:lw_friend>"`
}

func NewUser(name string, email string) (*User, error) {
	return &User{name: name, email: email}, nil
}

var errToManyWishes = errors.New("wish per user limit exceeded")

func (u *User) AddWish(text string) error {
	if u.wishes.Count() > 10 {
		return errToManyWishes
	}

	u.wishes.Add(newWish(text))

	return nil
}

type NotifyService interface {
	NotifyFriend(friend *User)
}

func (u *User) GrantWishes(notifier NotifyService) {
	wishesIt := u.wishes.MakeIter()
	for wishesIt.Next() {
		wishesIt.Value().(*Wish).grant()
	}

	friendsIt := u.friends.MakeIter()
	for friendsIt.Next() {
		notifier.NotifyFriend(friendsIt.Value().(*User))
	}
}

var ErrToManyFriends = errors.New("friend limit exceeded")

func (u *User) AddFriend(friend *User) error {
	if u.friends.Count() > 100 {
		return ErrToManyFriends
	}

	u.friends.Add(friend)
	return nil
}
