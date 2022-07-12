package telegram

import (
	"sync"
	"time"
)

type userUpdate struct {
	chatID          int64
	lastMessageTime int64
}

type usersUpdates struct {
	mu      sync.Mutex
	updates map[int]userUpdate
}

func newUsersUpdates() *usersUpdates {
	return &usersUpdates{
		mu:      sync.Mutex{},
		updates: make(map[int]userUpdate),
	}
}

func (u *usersUpdates) userCanReceiveMessage(userId int, t time.Time) bool {
	user, ok := u.updates[userId]
	if !ok {
		return true
	}

	return user.lastMessageTime+int64(time.Second) <= t.UnixNano()
}

func (u usersUpdates) add(userID int, update userUpdate) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.updates[userID] = update
}
