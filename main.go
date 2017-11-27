package main

import (
	"flag"
	"log"
	"strings"

	"net/http"

	"github.com/gorilla/websocket"
)

func validUser(user string) bool {
	return len(user) < 80 && !strings.HasPrefix(user, "#")
}

func makeWsHandler() func(http.ResponseWriter, *http.Request) {
	upgrader := websocket.Upgrader{}
	return func(w http.ResponseWriter, r *http.Request) {
		var usr string
		chatroom := r.URL.Query()["room"]
		username := r.URL.Query()["username"]
		roomname := "general"

		if len(chatroom) == 1 {
			roomname = chatroom[0]
		}

		if len(username) == 1 && validUser(username[0]) {
			usr = username[0]
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Need a valid username (shorter than 80 characters and does not start with a hash)"))
			return
		}

		ws, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Could not open websocket from connection."))
			return
		}

		room := joinOrCreateRoom(roomname, usr, ws)

		if room != nil {
			go handleMessages(usr, ws)
		}
	}
}

var (
	addr  string
	wsdir string
	dir   string
)

func init() {
	flag.StringVar(&addr,  "p",    ":8080",  "addr to use")
	flag.StringVar(&dir,   "root", "public", "http directory to serve")
	flag.StringVar(&wsdir, "ws",   "/ws",    "websocket handler path")
	flag.Parse()
}

func main() {

	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	http.HandleFunc(wsdir, makeWsHandler())

	err := http.ListenAndServe(addr, nil)

	if err != nil {
		log.Fatal("Server died with message:", err)
	}
}
