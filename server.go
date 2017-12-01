package main

import (
	"net"
	"log"
	"strings"
	"github.com/gorilla/websocket"
)

type Connection interface {
	WriteJSON(interface{}) error
	ReadJSON(interface{}) error
	RemoteAddr() net.Addr
	Close() error
}

type Server struct {
	Name  string
	Rooms map[string]*Room
	Users map[string]*Client
}

func NewServer(name string) *Server {
	return &Server{
		Name: name,
		Rooms: make(map[string]*Room),
		Users: make(map[string]*Client),
	}
}

func (s *Server) AddRoom(r *Room) error {
	name := r.Name
	if _, ok := s.Rooms[name]; ok {
		return errRoomExists
	}

	s.Rooms[name] = r
	return nil
}

func (s *Server) DelRoom(r *Room) error {
	name := r.Name
	if _, ok := s.Rooms[name]; !ok {
		return errNoSuchRoom
	}

	delete(s.Rooms, name)
	return nil
}

func (s *Server) GetRoom(n string) *Room {
	r, ok := s.Rooms[n]
	if !ok {
		r = NewRoom(n)
		s.Rooms[n] = r
	}

	return r
}

func (s *Server) AddUser(c *Client) error {
	name := c.Name
	if _, ok := s.Users[name]; ok {
		return errUserExists
	}

	s.Users[name] = c
	return nil
}

func (s *Server) DelUser(c *Client) error {
	name := c.Name
	if _, ok := s.Users[name]; !ok {
		return errNoSuchUser
	}

	delete(s.Users, name)
	return nil
}

func (s *Server) Nick(a, b string) error {
	if _, ok := s.Users[b]; ok {
		return errUserExists
	}

	if c, ok := s.Users[a]; ok {
		delete(s.Users, a)
		c.Name = b
		s.Users[b] = c
		return nil
	}

	return errNoSuchUser
}

func (s *Server) Compose(text string) Message {
	return Message{
		From: s.Name,
		Message: text,
		Meta: true,
	}
}

func (s *Server) RoomNames() []string {
	rs := make([]string, 0, len(s.Rooms))
	for n := range s.Rooms {
		rs = append(rs, n)
	}

	return rs
}

func (s *Server) Handle(c *Client) {
	ws := c.Conn
	defer ws.Close()

	for {
		var m Message
		err := ws.ReadJSON(&m)
		if err != nil {
			c.LeaveAll()
			s.DelUser(c)
			if isCloseError(err) {
				log.Printf("error: " + err.Error())
			}
			break
		}

		m.Meta = false
		m.From = c.Name

		r, ok := s.Rooms[m.Room]
		if !ok {
			sysMsg := s.Compose(m.Room + " room does not exist.")
			c.Write(sysMsg)
			continue
		}

		cmd := strings.SplitN(m.Message, " ", 2)
		if command, ok := commands[cmd[0]]; ok {
			command(s, r, c, cmd)
			continue
		}

		r.Write(m)
	}
}

func isCloseError(e error) bool {
	return !websocket.IsCloseError(e,
	  websocket.CloseNormalClosure,
	  websocket.CloseGoingAway,
	  websocket.CloseNoStatusReceived)
}
