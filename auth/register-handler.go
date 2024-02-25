package auth

import (
	"encoding/json"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func RegisterViewHandler(c echo.Context) error {
	register := register()

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return register.Render(c.Request().Context(), c.Response().Writer)
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)

		if err != nil || sess.Values["user"] == nil {
			c.Logger().Error("AuthMiddleware: session not found")
			return c.Redirect(302, "/register")
		}

		var user User
		err = json.Unmarshal(sess.Values["user"].([]byte), &user)
		if err != nil {
			c.Logger().Error("AuthMiddleware: user not found")
			return c.Redirect(302, "/register")
		}

		c.Set("user", user)
		return next(c)
	}
}

func RegisterHandler(c echo.Context) error {

	if err := c.Request().ParseForm(); err != nil {
		return err
	}

	var user User
	if err := c.Bind(&user); err != nil {
		c.Logger().Error("RegisterHandler: user not binded ", err)
		return c.Redirect(302, "/register")
	}

	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		c.Logger().Error("RegisterHandler: user not marshaled ", err)
		return err
	}

	sess.Values["user"] = userJson
	err = sess.Save(c.Request(), c.Response())

	UsersStore.Add(user)

	if err != nil {
		c.Logger().Error("RegisterHandler: session not saved ", err)
		return err
	}

	return c.Redirect(302, "/")
}
