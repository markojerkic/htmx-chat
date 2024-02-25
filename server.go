package main

import (
	"htmx-chat/auth"
	"htmx-chat/chat"

	"github.com/labstack/echo/v4"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
)

func main() {
	e := echo.New()
	e.Debug = true

	// e.Use(middleware.Logger())
	sessionMiddleware := session.Middleware(sessions.NewCookieStore([]byte("secret")))
	e.Use(sessionMiddleware)

	restricted := e.Group("/", sessionMiddleware, auth.AuthMiddleware);

	restricted.GET("/", chat.ChatHandler)
	restricted.GET("/room/:id", chat.RoomHandler)

	e.GET("/register", auth.RegisterViewHandler)
	e.POST("/register", auth.RegisterHandler)

	e.Static("/", "css")

	e.Logger.Fatal(e.Start(":1323"))
}
