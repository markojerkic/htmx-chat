package auth

import (
	"encoding/json"
	"fmt"
	"os"
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

func (r *RegisteredUsers) syncToFile() {
	// write r.users to byte array
	users, err := json.Marshal(r.users)
	if err != nil {
		fmt.Println("Error marshalling users")
		panic(err)
	}

	err = os.WriteFile("/tmp/chat-users", []byte(users), 0777)

	if err != nil {
		fmt.Println("Error writing users to file")
		panic(err)
	}
}

func newRegisteredUsers() *RegisteredUsers {
	users, err := os.ReadFile("/tmp/chat-users")

	var usersMap map[string]User
	if err != nil {
		fmt.Println("No users file found, creating new")
		usersMap = make(map[string]User)
	} else {
		err = json.Unmarshal(users, &usersMap)
		fmt.Println("Read users from file", usersMap)
		if err != nil {
			fmt.Println("Error unmarshalling users")
			panic(err)
		}
	}

	return &RegisteredUsers{
		users: usersMap,
		lock:  sync.RWMutex{},
	}
}

func (r *RegisteredUsers) Add(user User) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.users[user.ID] = user
	r.syncToFile()
}

// Get list of all users which are not the user querying
func (r *RegisteredUsers) GetAllUsers(c *echo.Context) []User {
	r.lock.RLock()
	defer r.lock.RUnlock()

	currentUsrer := (*c).Get("user").(User)

	var users = make([]User, 0, len(r.users))

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
