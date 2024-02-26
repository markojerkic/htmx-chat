package room

import (
	"htmx-chat/auth"
	"strings"

	"github.com/labstack/echo/v4"
)

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

	return newRoom(users).Render(c.Request().Context(), c.Response().Writer)
}

func AllRoomsHandler(c echo.Context) error {
	currentUsrer := c.Get("user").(auth.User)

	rooms := RoomsStore.GetAllMyRooms(currentUsrer)

	roomsComponent := allRooms(rooms, currentUsrer)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return roomsComponent.Render(c.Request().Context(), c.Response().Writer)
}

func OpenChatHandler(c echo.Context) error {
	currentUsrer := c.Get("user").(auth.User)

	requestUserId := c.Param("id")

	requestUser := auth.UsersStore.GetUserById(requestUserId)
	if requestUser == nil {
		return c.String(404, "User not found")
	}

	room := RoomsStore.AddRoom(currentUsrer, *requestUser)

	c.Logger().Infof("Room: %v", room)

	rooms := RoomsStore.GetAllMyRooms(currentUsrer)

	roomsComponent := createNewRoom(rooms, requestUserId, currentUsrer)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return roomsComponent.Render(c.Request().Context(), c.Response().Writer)
}
