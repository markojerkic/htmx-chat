package chat

import (
	"htmx-chat/models"
	"htmx-chat/templates"
	"os"

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

	roomsComponent.Render(c.Request().Context(), os.Stdout)
	c.Logger().Info()

	return roomsComponent.Render(c.Request().Context(), c.Response().Writer)
}
