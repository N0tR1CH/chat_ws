package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

type room struct {
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

func handleRoomWs(logger Logger) http.Handler {
	room := newRoom()
	go room.run()
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			socket, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				logger.Error("websocket connection was not possible", "err", err.Error())
				return
			}
			client := &client{
				socket: socket,
				send:   make(chan []byte, socketBufferSize),
				room:   room,
			}
			room.join <- client
			defer func() {
				room.leave <- client
			}()
			go client.write()
			client.read()
		},
	)
}
