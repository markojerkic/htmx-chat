package auth

import (
	"htmx-chat/db"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID   string `form:"username" json:"username"`
	Name string `form:"name" json:"name"`
}

func (u User) GetID() string {
	return u.ID
}

func (u User) SetID(ID string) db.Item {
	u.ID = ID
	return u
}

type RegisteredUsers struct {
	users *db.Collection[User]
}

func newRegisteredUsers() *RegisteredUsers {
	return &RegisteredUsers{
		users: db.NewCollection[User]("users"),
	}
}

func (r *RegisteredUsers) Add(user User) {
	r.users.Save(user)
}

func (r *RegisteredUsers) GetUserById(id string) (*User, error) {
	return r.users.Get(id)
}

// Get list of all users which are not the user querying
func (r *RegisteredUsers) GetAllUsers(c *echo.Context) []User {
	currentUsrer := (*c).Get("user").(User)
	return r.users.GetAllByPredicate(func(user User) bool {
		return user.ID != currentUsrer.ID
	})
}

var UsersStore = newRegisteredUsers()
