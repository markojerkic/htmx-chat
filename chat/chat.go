package chat

import (
	"htmx-chat/models"
	"htmx-chat/templates"

	"github.com/labstack/echo/v4"
)

func ChatHandler(c echo.Context) error {

	rooms := []models.ChatRoom{
		{
			ID:      "1",
			ClientA: &models.Client{ID: "1"},
			ClientB: &models.Client{ID: "2"},
		},
		{
			ID:      "2",
			ClientA: &models.Client{ID: "1"},
			ClientB: &models.Client{ID: "2"},
		},
	}

	roomsComponent := templates.Chat(rooms)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return roomsComponent.Render(c.Request().Context(), c.Response().Writer)
}

func RoomHandler(c echo.Context) error {

	c.Logger().Info("Room {}", c.Param("id"))
	roomComponent := templates.Room()

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return roomComponent.Render(c.Request().Context(), c.Response().Writer)
}
