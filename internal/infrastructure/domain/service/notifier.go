package service

import (
	"log"
	"lw/internal/domain"
)

type NotifyService struct {
}

func (n *NotifyService) NotifyFriend(friend *domain.User) {
	log.Printf("user %s notified", friend.Email())
}
