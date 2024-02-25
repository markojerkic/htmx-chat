package room

import (
	"htmx-chat/auth"
	"htmx-chat/templates"

	"github.com/labstack/echo/v4"
)

type client struct {
	ID string
}

type chatRoom struct {
	ID      string
	ClientA *client
	ClientB *client
}

func AllRoomsHandler(c echo.Context) error {
	currentUsrer := c.Get("user").(auth.User)

	rooms := []chatRoom{
		{
			ID:      "1",
			ClientA: &client{ID: "1"},
			ClientB: &client{ID: "2"},
		},
		{
			ID:      "2",
			ClientA: &client{ID: "1"},
			ClientB: &client{ID: "2"},
		},
	}

	roomsComponent := allRooms(rooms, currentUsrer)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return roomsComponent.Render(c.Request().Context(), c.Response().Writer)
}

func RoomHandler(c echo.Context) error {
	roomComponent := templates.Room()

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return roomComponent.Render(c.Request().Context(), c.Response().Writer)
}
