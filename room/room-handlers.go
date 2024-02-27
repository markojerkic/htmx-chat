package room

import (
	"fmt"
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

func SearchUsersFormHandler(c echo.Context) error {

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	users := auth.UsersStore.GetAllUsers(&c)

	return newRoom(users).Render(c.Request().Context(), c.Response().Writer)
}

func AllRoomsHandler(c echo.Context) error {
	currentUsrer := c.Get("user").(auth.User)

	rooms := RoomsStore.GetAllMyRooms(currentUsrer)

	var selectedRoom *chatRoom

	if c.Param("id") != "" {
		c.Logger().Debug("Request for a room")
		room, err := RoomsStore.GetRoom(c.Param("id"))

		if err != nil {
			return c.Redirect(302, "/")
		}
		if !room.IsUserInRoom(c.Get("user").(auth.User).ID) {
			return c.Redirect(302, "/")
		}

		selectedRoom = &room
	}

	roomsComponent := allRooms(rooms, currentUsrer, selectedRoom)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return roomsComponent.Render(c.Request().Context(), c.Response().Writer)
}

func createRoomWithRoomListPartialHandler(c echo.Context) error {
	currentUsrer := c.Get("user").(auth.User)

	requestUserId := c.Param("id")

	requestUser, err := auth.UsersStore.GetUserById(requestUserId)
	if err != nil {
		return c.String(404, "User not found")
	}

	room, err := RoomsStore.AddRoom(currentUsrer, requestUser)
	if err != nil {
		return c.String(500, "Error creating room")
	}
	if !room.IsUserInRoom(c.Get("user").(auth.User).ID) {
		c.Response().Header().Set("HX-Redirect", "/")
		return c.String(403, "You are not allowed to see this room")
	}

	c.Logger().Infof("Room: %v", room)

	rooms := RoomsStore.GetAllMyRooms(currentUsrer)

	roomsComponent := createNewRoom(rooms, room, currentUsrer)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	c.Logger().Infof("Room: %s", room.ID)
	c.Response().Header().Set("Hx-Push-Url", fmt.Sprintf("/room/%s", room.ID))

	return roomsComponent.Render(c.Request().Context(), c.Response().Writer)
}

// Check if the request is an HX request
// If it is, get the room ID from the request and get the room from the store
func RoomHandler(c echo.Context) error {
	if c.Request().Header.Get("Hx-Request") == "true" {
		if c.Request().Header.Get("Hx-Target") == "chat" {
			// Clicked on the room from the list
			return openRoomPartialHandler(c)
		}
		// Creating new room from search result
		return createRoomWithRoomListPartialHandler(c)
	}

	// Normal request, return the full page
	return AllRoomsHandler(c)
}
