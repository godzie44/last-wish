package service

import (
	"github.com/godzie44/lw/internal/domain"
	"log"
)

type NotifyService struct {
}

func (n *NotifyService) NotifyFriend(friend *domain.User) {
	log.Printf("user %s notified", friend.Email())
}
