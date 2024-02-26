package room

import (
	"encoding/json"
	"fmt"
	"htmx-chat/auth"
	"os"
	"sync"

	guuid "github.com/google/uuid"
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

type roomStore struct {
	lock  sync.RWMutex
	rooms map[string]chatRoom
}

func (r *roomStore) GetAllMyRooms(user auth.User) []chatRoom {
	r.lock.RLock()
	defer r.lock.RUnlock()

	myRooms := make([]chatRoom, 0, 10)
	for _, room := range r.rooms {
		if room.ClientA.ID == user.ID || room.ClientB.ID == user.ID {
			myRooms = append(myRooms, room)
		}
	}

	return myRooms
}

func (r *roomStore) AddRoom(userA auth.User, userB auth.User) chatRoom {
	r.lock.Lock()
	defer r.lock.Unlock()

	for id, room := range r.rooms {
		isClientA := room.ClientA.ID == userA.ID || room.ClientA.ID == userB.ID
		isClientB := room.ClientB.ID == userA.ID || room.ClientB.ID == userB.ID

		if isClientA || isClientB {
			return r.rooms[id]
		}
	}

	room := chatRoom{
		ID:      guuid.NewString(),
		ClientA: &userA,
		ClientB: &userB,
	}

	r.rooms[room.ID] = room

	r.syncToFile()

	return room
}

func (r *roomStore) GetRoom(id string) chatRoom {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.rooms[id]
}

func (r *roomStore) syncToFile() {
	// write r.rooms to byte array
	rooms, err := json.Marshal(r.rooms)
	if err != nil {
		fmt.Println("Error marshalling rooms")
		panic(err)
	}

	err = os.WriteFile("/tmp/chat-rooms", []byte(rooms), 0777)

	if err != nil {
		fmt.Println("Error writing users to file")
		panic(err)
	}
}

func newRoomStore() *roomStore {
	rooms, err := os.ReadFile("/tmp/chat-rooms")

	var roomsMap map[string]chatRoom
	if err != nil {
		fmt.Println("No rooms file found, creating new")
		roomsMap = make(map[string]chatRoom)
	} else {
		err = json.Unmarshal(rooms, &roomsMap)
		fmt.Println("Read rooms from file", roomsMap)
		if err != nil {
			fmt.Println("Error unmarshalling rooms")
			panic(err)
		}
	}

	return &roomStore{
		rooms: roomsMap,
		lock:  sync.RWMutex{},
	}
}

var RoomsStore = newRoomStore()
