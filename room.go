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
	done    chan struct{}
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
		case <-r.done:
			close(r.forward)
			close(r.join)
			close(r.leave)
			close(r.done)
		}
	}
}

func handleRoomWs(logger Logger) http.Handler {
	rooms := make(map[string]*room)
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			num := r.URL.Query().Get("num")
			if _, ok := rooms[num]; !ok {
				rooms[num] = newRoom()
				go rooms[num].run()
			}

			socket, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				logger.Error("websocket connection was not possible", "err", err.Error())
				return
			}
			client := &client{
				socket: socket,
				send:   make(chan []byte, socketBufferSize),
				room:   rooms[num],
			}
			rooms[num].join <- client
			defer func() {
				rooms[num].leave <- client
				logger.Info("handleRoomWs - Closing Room", "roomNumber", num)
				// Below must be a channel
				if len(rooms[num].clients) == 0 {
					rooms[num].done <- struct{}{}
					logger.Info("handleRoomWs - Closing Room", "roomNumber", num)
				}
			}()
			go client.write()
			client.read()
		},
	)
}
