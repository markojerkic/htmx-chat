package room

import (
	"htmx-chat/auth"
	"htmx-chat/db"

	"github.com/gorilla/websocket"
)

type chatRoom struct {
	ID        string
	ClientAID string
	ClientBID string
	wsClients map[string]*websocket.Conn `json:'-'`
}

func (c chatRoom) GetUserA() (*auth.User, error) {
	return auth.UsersStore.GetUserById(c.ClientAID)
}

func (c chatRoom) GetUserB() (*auth.User, error) {
	return auth.UsersStore.GetUserById(c.ClientBID)
}

func (r *chatRoom) GetClientWhichIsNotMe(myId string) auth.User {
	if r.ClientAID == myId {
		user, err := r.GetUserB()
		if err != nil {
			panic(err)
		}
		return *user
	}
	user, err := r.GetUserA()
	if err != nil {
		panic(err)
	}
	return *user
}

func (r *chatRoom) IsUserInRoom(userID string) bool {
	return r.ClientAID == userID || r.ClientBID == userID
}

func (r chatRoom) GetID() string {
	return r.ID
}

func (r chatRoom) SetID(ID string) db.Item {
	r.ID = ID
	return r
}

type roomStore struct {
	rooms *db.Collection[chatRoom]
}

func (r *roomStore) GetAllMyRooms(user auth.User) []chatRoom {
	return r.rooms.GetAllByPredicate(func(room chatRoom) bool {
		return room.ClientAID == user.ID || room.ClientBID == user.ID
	})
}

func (r *roomStore) AddRoom(userA auth.User, userB auth.User) (chatRoom, error) {
	room, err := r.rooms.GetByPredicate(func(room chatRoom) bool {
		isClientA := room.ClientAID == userA.ID || room.ClientAID == userB.ID
		isClientB := room.ClientBID == userA.ID || room.ClientBID == userB.ID

		return isClientA && isClientB
	})

	if err == nil {
		return room, nil
	}

	room = chatRoom{
		ClientAID: userA.ID,
		ClientBID: userB.ID,
		wsClients: make(map[string]*websocket.Conn),
	}

	return r.rooms.Save(room)
}

func (r *roomStore) GetRoom(id string) (*chatRoom, error) {
	room, err := r.rooms.Get(id)

	if err == nil && room.wsClients == nil {
		room.wsClients = make(map[string]*websocket.Conn)
	}

	return room, err
}

func newRoomStore() *roomStore {
	return &roomStore{
		rooms: db.NewCollection[chatRoom]("rooms"),
	}
}

var RoomsStore = newRoomStore()
