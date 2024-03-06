package room

import (
	"bytes"
	"context"
	"encoding/json"
	"htmx-chat/auth"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Message  string    `json:"message"`
	RoomId   string    `json:"roomId"`
	SenderId string    `json:"senderId"`
	Date     time.Time `json:"date"`
}

type Client struct {
	wsConnection    *websocket.Conn
	userId          string
	messageReceiver chan Message
	hub             *Hub
	logger          echo.Logger
	currentRoomId   string
}

type WsMessage struct {
	RoomId  string            `json:"roomId"`
	Message string            `json:"message"`
	HEADERS map[string]string `json:"HEADERS"`
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

		if receivedWsMessage.HEADERS["HX-Trigger"] == "currentRoom" {
			c.currentRoomId = receivedWsMessage.RoomId
			continue
		}

		message := Message{
			Message:  receivedWsMessage.Message,
			RoomId:   receivedWsMessage.RoomId,
			SenderId: c.userId,
			Date:     time.Now(),
		}

		room, err := RoomsStore.rooms.Get(receivedWsMessage.RoomId)
		if err != nil {
			c.logger.Errorf("Error: %v", err)
			continue
		}

		room.AddMessage(message)
		c.hub.messages <- message
	}
}

func (c *Client) sendMessages() {
	defer c.wsConnection.Close()
	for {
		message := <-c.messageReceiver

		renderedMessage := new(bytes.Buffer)

		if c.currentRoomId == message.RoomId {
			chatBubble(false, message.Message, message.Date).Render(context.Background(), renderedMessage)

			err := c.wsConnection.WriteMessage(websocket.TextMessage, renderedMessage.Bytes())
			if err != nil {
				c.logger.Errorf("Error: %v", err)
				break
			}
            renderedMessage.Reset()
		}

		myRooms := RoomsStore.GetAllMyRooms(auth.User{ID: c.userId})
		roomListWithPreview(myRooms, c.userId).Render(context.Background(), renderedMessage)

		err := c.wsConnection.WriteMessage(websocket.TextMessage, renderedMessage.Bytes())
		if err != nil {
			c.logger.Errorf("Error: %v", err)
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
