package room

import (
	"bytes"
	"context"
	"fmt"
	"htmx-chat/auth"

	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

func openRoomPartialHandler(c echo.Context) error {
	roomId := c.Param("id")

	room, err := RoomsStore.GetRoom(roomId)
	if err != nil {
		c.Logger().Error("Error getting requested room", err)
		return c.String(404, "Room not found")
	}
	if !room.IsUserInRoom(c.Get("user").(auth.User).ID) {
		c.Response().Header().Set("HX-Redirect", "/")
		return c.String(403, "You are not allowed to see this room")
	}

	c.Logger().Debugf("Room: %v", room)

	currentUser := c.Get("user").(auth.User)
	requestedUser := room.GetClientWhichIsNotMe(currentUser.ID)

	c.Logger().Debugf("Requested user: %v", requestedUser)

	chatComponent := Chat(requestedUser.ID, requestedUser.Name, room.ID)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return chatComponent.Render(c.Request().Context(), c.Response().Writer)
}

func ConnectToRoom(c echo.Context) error {
	c.Logger().Error("Want to connect to a room")
	roomId := c.Param("id")
	room, err := RoomsStore.GetRoom(roomId)

	if err != nil {
		c.Logger().Errorf("Room %s does not exist", roomId)
		return c.String(404, fmt.Sprintf("%v", err))
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		c.Logger().Error("Error creating ws connection", err)
		return err
	}
	defer ws.Close()

	room.wsClients[roomId] = ws

	for {

		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
			return err
		}
		var message map[string]string
		err = json.Unmarshal(msg, &message)
		if err != nil {
			c.Logger().Error(err)
			return err
		}

		renderedMessage := new(bytes.Buffer)
		chatBubble(true, message["message"]).Render(context.Background(), renderedMessage)

		ws.WriteMessage(websocket.TextMessage, renderedMessage.Bytes())
	}
}
