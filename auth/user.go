package auth

import "sync"

type User struct {
	ID   string `form:"id" json:"id"`
	Name string `form:"name" json:"name"`
}

type RegisteredUsers struct {
	lock  sync.RWMutex
	users map[string]User
}

func newRegisteredUsers() *RegisteredUsers {
	return &RegisteredUsers{
		users: make(map[string]User),
		lock:  sync.RWMutex{},
	}
}

func (r *RegisteredUsers) Add(user User) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.users[user.ID] = user
}

var UsersStore = newRegisteredUsers()
