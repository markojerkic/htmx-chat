package main

import (
	"htmx-chat/auth"
	"htmx-chat/room"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func main() {
	hub := room.NewHub()
    go hub.Run()

	e := echo.New()
	e.Debug = true

	sessionMiddleware := session.Middleware(sessions.NewCookieStore([]byte("secret")))
	e.Use(sessionMiddleware)

	restricted := e.Group("", auth.AuthMiddleware)

	restricted.GET("/", room.AllRoomsHandler)

	restricted.GET("/room/:id", room.RoomHandler)

	restricted.GET("/room/new", room.SearchUsersFormHandler)
	restricted.GET("/room/search", room.SearchUsersNewRoom)

	restricted.GET("/ws", room.ConnectClientToWS(hub))

	e.GET("/register", auth.RegisterViewHandler)
	e.POST("/register", auth.RegisterHandler)

	e.Static("/", "css")

	e.Logger.Fatal(e.Start(":1323"))
}
