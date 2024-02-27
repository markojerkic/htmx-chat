package room

import (
	"htmx-chat/auth"
	"htmx-chat/db"
)

type chatRoom struct {
	ID      string
	ClientA *auth.User
	ClientB *auth.User
}

func (r *chatRoom) GetClientWhichIsNotMe(myId string) auth.User {
	if r.ClientA.ID == myId {
		return auth.User(*r.ClientB)
	}
	return auth.User(*r.ClientA)
}

func (r *chatRoom) IsUserInRoom(userID string) bool {
	return r.ClientA.ID == userID || r.ClientB.ID == userID
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
		return room.ClientA.ID == user.ID || room.ClientB.ID == user.ID
	})
}

func (r *roomStore) AddRoom(userA auth.User, userB auth.User) (chatRoom, error) {
	room, err := r.rooms.GetByPredicate(func(room chatRoom) bool {
		isClientA := room.ClientA.ID == userA.ID || room.ClientA.ID == userB.ID
		isClientB := room.ClientB.ID == userA.ID || room.ClientB.ID == userB.ID

		return isClientA && isClientB
	})

	if err == nil {
		return room, nil
	}

	room = chatRoom{
		ClientA: &userA,
		ClientB: &userB,
	}

	return r.rooms.Save(room)
}

func (r *roomStore) GetRoom(id string) (chatRoom, error) {
	return r.rooms.Get(id)
}

func newRoomStore() *roomStore {
	return &roomStore{
		rooms: db.NewCollection[chatRoom]("rooms"),
	}
}

var RoomsStore = newRoomStore()
