package main

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

type wsConnClient struct {
	*websocket.Conn
}

func (c *wsConnClient) Send(m *Message) {
	c.WriteJSON(m)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	wsc := &wsConnClient{ws}
	p := r.URL.Path
	router.AddRoute(p, wsc)
	defer router.DelRoute(p, wsc)

	for {
		m := &Message{}
		if err := ws.ReadJSON(m); err != nil {
			if isCloseError(err) {
				log.Printf("error: " + err.Error())
			}
			break
		}
		router.Send(m)
	}
}

func isCloseError(e error) bool {
	return !websocket.IsCloseError(e,
		websocket.CloseNormalClosure,
		websocket.CloseGoingAway,
		websocket.CloseNoStatusReceived)
}
