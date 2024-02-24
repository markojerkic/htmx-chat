package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Marko Jerkić!")
	})
    e.Logger.Info("Server started at :1323")
	e.Logger.Fatal(e.Start(":1323"))
}
