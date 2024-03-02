package room

import (
	"htmx-chat/auth"
	"htmx-chat/db"
)

type chatRoom struct {
	ID        string
	ClientAID string
	ClientBID string
	Messages  []Message
}

func (c chatRoom) GetUserA() (*auth.User, error) {
	return auth.UsersStore.GetUserById(c.ClientAID)
}

func (c chatRoom) GetUserB() (*auth.User, error) {
	return auth.UsersStore.GetUserById(c.ClientBID)
}

func (c *chatRoom) AddMessage(message Message) {
	c.Messages = append(c.Messages, message)
	RoomsStore.Save(*c)
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

func (r *roomStore) Save(room chatRoom) (chatRoom, error) {
	return r.rooms.Save(room)
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
		Messages:  make([]Message, 10),
	}

	return r.rooms.Save(room)
}

func (r *roomStore) GetRoom(id string) (*chatRoom, error) {
	room, err := r.rooms.Get(id)

	return room, err
}

func newRoomStore() *roomStore {
	return &roomStore{
		rooms: db.NewCollection[chatRoom]("rooms"),
	}
}

var RoomsStore = newRoomStore()
