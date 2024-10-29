package main

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

type clientRes struct {
	Msg     string         `json:"msg"`
	Val     string         `json:"val"`
	Headers map[string]any `json:"HEADERS"`
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		slog.Info("client#write", "msg", string(msg))
		html := `<div hx-swap-oob="beforeend:#notifications"><p>%s</p></div>`
		res := new(clientRes)
		json.Unmarshal(msg, res)
		slog.Info("client#write", "m", res)
		msg = []byte(fmt.Sprintf(html, string(res.Msg)))
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}
