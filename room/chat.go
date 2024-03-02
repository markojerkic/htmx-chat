package room

import (
	"encoding/json"
	"htmx-chat/auth"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	message  string
	roomId   string
	senderId string
}

type Client struct {
	wsConnection    *websocket.Conn
	userId          string
	messageReceiver chan Message
	hub             *Hub
	logger          echo.Logger
}

type WsMessage struct {
	Message string            `json:message`
	RoomId  string            `json:roomId`
	HEADERS map[string]string `json:HEADERS`
}

func (c *Client) readMessages() {
	defer c.wsConnection.Close()

	var receivedWsMessage WsMessage
	for {
		_, messageBuf, err := c.wsConnection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Errorf("Error: %v", err)
			}
			break
		}

		json.Unmarshal(messageBuf, &receivedWsMessage)

		c.logger.Debug("Received message: ", receivedWsMessage.Message)

		message := Message{
			message:  receivedWsMessage.Message,
			roomId:   receivedWsMessage.RoomId,
			senderId: c.userId,
		}

		c.hub.messages <- message
	}
}

func (c *Client) sendMessages() {
	defer c.wsConnection.Close()
	for {
		message := <-c.messageReceiver

		err := c.wsConnection.WriteJSON(message)
		if err != nil {
			break
		}
	}

}

func ConnectClientToWS(hub *Hub) func(echo.Context) error {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(auth.User)
		if !ok {
			return c.Redirect(302, "/register")
		}

		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}

		client := &Client{
			wsConnection:    ws,
			hub:             hub,
			userId:          user.ID,
			messageReceiver: make(chan Message, 256),
			logger:          c.Logger(),
		}

		hub.connectionRequest <- client

		go client.readMessages()
		go client.sendMessages()

		return nil
	}

}
