package chat

import (
	"github.com/labstack/echo/v4"
)

func OpenChatHandler(c echo.Context) error {

	chatComponent := Chat()

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return chatComponent.Render(c.Request().Context(), c.Response().Writer)
}

