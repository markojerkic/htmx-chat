package main

import (
	"htmx-chat/auth"
	"htmx-chat/room"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Debug = true

	sessionMiddleware := session.Middleware(sessions.NewCookieStore([]byte("secret")))
	e.Use(sessionMiddleware)

	restricted := e.Group("", auth.AuthMiddleware);

	restricted.GET("/", room.AllRoomsHandler)

	restricted.GET("/room/:id/open", room.OpenRoomPartialHandler)
	restricted.GET("/room/:id", room.OpenRoomHandler)

	restricted.GET("/room/new", room.SearchUsersFormHandler)
	restricted.GET("/room/search", room.SearchUsersNewRoom)

	e.GET("/register", auth.RegisterViewHandler)
	e.POST("/register", auth.RegisterHandler)

	e.Static("/", "css")

	e.Logger.Fatal(e.Start(":1323"))
}
