package main

import (
	"htmx-chat/templates"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		helloTemplate := templates.Hello("Marko")
		c.Logger().Info("HelloTemplate: ", helloTemplate)
		return helloTemplate.Render(c.Request().Context(), c.Response().Writer)
	})

	e.Logger.Info("Server started at :1323")
	e.Logger.Fatal(e.Start(":1323"))
}
