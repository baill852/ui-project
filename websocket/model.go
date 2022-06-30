package ws

import "github.com/gorilla/websocket"

type Server interface {
	AddClient(client Client)
	RemoveClient(client Client)
	Publish(data interface{})
}

type Client struct {
	Id   string
	Conn *websocket.Conn
}
