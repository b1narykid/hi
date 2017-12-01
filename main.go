package main

import (
	"log"
	"flag"

	"net/http"

	"github.com/gorilla/websocket"
)

func httpHandler(w http.ResponseWriter, r *http.Request, s *Server) {
	upgrader := websocket.Upgrader{}
	q := r.URL.Query()

	name := q.Get("username")
	if !UserNameIsValid(name) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errErroneusNickname.Error()))
		return
	}

	c := NewClient(name, nil)
	c.Name = name

	err := s.AddUser(c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errUserExists.Error()))
		return
	}

	if c.Conn, err = upgrader.Upgrade(w, r, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errOpenWsConn.Error()))
		return
	}

	room := s.GetRoom("#general")
	c.Join(room)

	s.HandleClient(c)
}

var (
	addr  string
	dir   string
	wsdir string
	fsdir string
)

func init() {
	flag.StringVar(&addr,  "bind", ":8080",  "bind to this addr")
	flag.StringVar(&dir,   "root", "public", "http directory to serve")
	flag.StringVar(&wsdir, "ws",   "/ws",    "websocket handler path")
	flag.StringVar(&fsdir, "fs",   "/",      "file server handler path")
	flag.Parse()
}

func main() {
	if fsdir != "" {
		fs := http.FileServer(http.Dir(dir))
		http.Handle("/", fs)
	}

	s := NewServer("hi")

	http.HandleFunc(wsdir, func(w http.ResponseWriter, r *http.Request) {
		httpHandler(w, r, s)
	})

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("Server died with message: ", err)
	}
}
