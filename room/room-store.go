package room

import (
	"fmt"
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

func (r chatRoom) GetID() string {
	return r.ID
}

func (r chatRoom) SetID(ID string) {
	r.ID = ID
}

type roomStore struct {
	rooms *db.Collection[chatRoom]
}

func (r *roomStore) GetAllMyRooms(user auth.User) []chatRoom {
	rooms, err := r.rooms.GetAll()
	if err != nil {
		fmt.Println("Error getting all rooms")
		panic(err)
	}
	myRooms := make([]chatRoom, 0, len(rooms))

	for _, room := range rooms {
		if room.ClientA.ID == user.ID || room.ClientB.ID == user.ID {
			myRooms = append(myRooms, room)
		}
	}

	return myRooms
}

func (r *roomStore) AddRoom(userA auth.User, userB auth.User) chatRoom {
	rooms, err := r.rooms.GetAll()
	if err != nil {
		fmt.Println("Error getting all rooms")
		panic(err)
	}

	for id, room := range rooms {
		isClientA := room.ClientA.ID == userA.ID || room.ClientA.ID == userB.ID
		isClientB := room.ClientB.ID == userA.ID || room.ClientB.ID == userB.ID

		if isClientA && isClientB {
			return rooms[id]
		}
	}

	room := chatRoom{
		ClientA: &userA,
		ClientB: &userB,
	}

	r.rooms.Save(room)
	return room
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
