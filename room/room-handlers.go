package room

import (
	"htmx-chat/auth"
	"strings"

	"github.com/labstack/echo/v4"
)

type client auth.User

type chatRoom struct {
	ID      string
	ClientA *client
	ClientB *client
}

func (r *chatRoom) GetClientWhichIsNotMe(myId string) auth.User {
	if r.ClientA.ID == myId {
		return auth.User(*r.ClientB)
	}

	return auth.User(*r.ClientA)
}

func SearchUsersNewRoom(c echo.Context) error {
	users := auth.UsersStore.GetAllUsers(&c)

	filteredUsers := make([]auth.User, 0)

	for _, user := range users {
		query := strings.ToLower(c.QueryParam("q"))

		idContains := strings.Contains(strings.ToLower(user.ID), query)
		nameContains := strings.Contains(strings.ToLower(user.Name), query)

		if idContains || nameContains {
			filteredUsers = append(filteredUsers, user)
		}

	}

	return newRoomUsersList(filteredUsers).Render(c.Request().Context(), c.Response().Writer)
}

func NewRoomHandler(c echo.Context) error {

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	users := auth.UsersStore.GetAllUsers(&c)
	c.Logger().Infof("Users: %v", users)

	return newRoom(users).Render(c.Request().Context(), c.Response().Writer)
}

func AllRoomsHandler(c echo.Context) error {
	currentUsrer := c.Get("user").(auth.User)

	rooms := []chatRoom{}

	roomsComponent := allRooms(rooms, currentUsrer)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return roomsComponent.Render(c.Request().Context(), c.Response().Writer)
}

func OpenChatHandler(c echo.Context) error {
	currentUsrer := c.Get("user").(auth.User)

	rooms := []chatRoom{}

	requestUserId := c.Param("id")

	roomsComponent := createNewRoom(rooms, requestUserId, currentUsrer)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return roomsComponent.Render(c.Request().Context(), c.Response().Writer)
}
