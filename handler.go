package main

import (
	"net/http"
	"github.com/gorilla/websocket"
)

type wsClient struct {
	*websocket.Conn
	Path string
}

func (c *wsClient) Send(m *Message) {
	c.WriteJSON(m)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	go wsServe(router, &wsClient{
		Conn: ws,
		Path: r.URL.Path,
	})
}

func wsServe(r *Router, c *wsClient) {
	r.AddRoute(c.Path, c)
	for {
		m := &Message{}
		if err := c.ReadJSON(m); err != nil {
			break
		}
		r.Send(m)
	}
	r.DelRoute(c.Path, c)
	c.Close()
}
