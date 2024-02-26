package room

import (
	"htmx-chat/auth"

	"github.com/labstack/echo/v4"
)

func OpenRoomPartialHandler(c echo.Context) error {

	roomId := c.Param("id")

	room, err := RoomsStore.GetRoom(roomId)

	if err != nil {
		c.Logger().Error("Error getting requested room", err)
		return c.String(404, "Room not found")
	}

	currentUser := c.Get("user").(auth.User)
	requestedUser := room.GetClientWhichIsNotMe(currentUser.ID)

	chatComponent := Chat(requestedUser.ID, requestedUser.Name)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return chatComponent.Render(c.Request().Context(), c.Response().Writer)
}
