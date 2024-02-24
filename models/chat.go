package models

type Client struct {
	ID string
}

type ChatRoom struct {
	ID      string
	ClientA *Client
	ClientB *Client
}
