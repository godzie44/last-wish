package domain

import (
	"database/sql"
	"errors"
	"github.com/godzie44/d3/orm/entity"
	"unicode/utf8"
)

type Email string

func (e Email) isValid() bool {
	return utf8.RuneCountInString(string(e)) < 255
}

//d3:entity
//d3_table:lw_user
type User struct {
	id      sql.NullInt64      `d3:"pk:auto"`
	name    string             `d3:"column:name"`
	email   Email              `d3:"column:email"`
	wishes  *entity.Collection `d3:"one_to_many:<target_entity:lw/internal/domain/Wish,join_on:user_id,delete:nullable>"`
	friends *entity.Collection `d3:"many_to_many:<target_entity:lw/internal/domain/User,join_on:u1_id,reference_on:u2_id,join_table:lw_friend>"`
}

const maxWishesPerUser = 10
const maxFriendsPerUser = 100

var errToManyWishes = errors.New("wish per user limit exceeded")
var errInvalidEmail = errors.New("email is invalid")
var errUserNameTooLong = errors.New("name is too long")

func NewUser(name string, email string) (*User, error) {
	em := Email(email)
	if !em.isValid() {
		return nil, errInvalidEmail
	}

	if utf8.RuneCountInString(name) > 40 {
		return nil, errUserNameTooLong
	}

	return &User{name: name, email: em, wishes: entity.NewCollection(), friends: entity.NewCollection()}, nil
}

func (u *User) MakeWish(text string) error {
	if u.wishes.Count() >= maxWishesPerUser {
		return errToManyWishes
	}

	wish, err := newWish(text)
	if err != nil {
		return err
	}

	u.wishes.Add(wish)

	return nil
}

type NotifyService interface {
	NotifyFriend(friend *User)
}

func (u *User) ReleaseWishes(notifier NotifyService) {
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
	if u.friends.Count() >= maxFriendsPerUser {
		return ErrToManyFriends
	}

	u.friends.Add(friend)
	return nil
}

func (u *User) Email() string {
	return string(u.email)
}

func (u *User) FriendsCount() int {
	return u.friends.Count()
}
