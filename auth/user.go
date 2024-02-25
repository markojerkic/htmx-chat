package auth

import (
	"sync"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID   string `form:"username" json:"username"`
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

// Get list of all users which are not the user querying
func (r *RegisteredUsers) GetAllUsers(c *echo.Context) []User {
	r.lock.RLock()
	defer r.lock.RUnlock()

	currentUsrer := (*c).Get("user").(User)

	var users = make([]User, len(r.users))

	for _, user := range r.users {
		if user.ID != currentUsrer.ID {
			users = append(users, user)
		} else {
			(*c).Logger().Infof("Skipping user %s", user.ID)
		}
	}
	return users
}

var UsersStore = newRegisteredUsers()
