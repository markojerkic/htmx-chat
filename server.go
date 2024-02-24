package main

import (
	"htmx-chat/chat"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", chat.ChatHandler)

	e.Static("/", "css")

	e.Logger.Info("Server started at :1323")
	e.Logger.Fatal(e.Start(":1323"))
}
