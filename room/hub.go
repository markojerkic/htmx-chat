package room

type Hub struct {
	connectionRequest chan *Client
	disconnectClient  chan *Client
	messages          chan Message

	connectedClients map[*Client]bool
}

func (h *Hub) broadcastMessage(message Message, receiverId string) {
	for client := range h.connectedClients {
		if client.userId == receiverId {
			client.messageReceiver <- message
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.connectionRequest:
			h.connectedClients[client] = true
		case client := <-h.disconnectClient:
			delete(h.connectedClients, client)
			close(client.messageReceiver)
		case message := <-h.messages:
			roomId := message.roomId
			room, err := RoomsStore.GetRoom(roomId)
			if err != nil {
				continue
			}

			receiver := room.GetClientWhichIsNotMe(message.senderId)
			h.broadcastMessage(message, receiver.ID)
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		connectionRequest: make(chan *Client, 256),
		disconnectClient:  make(chan *Client, 256),
		messages:          make(chan Message, 256),
		connectedClients:  make(map[*Client]bool, 256),
	}
}
