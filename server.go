package main

import (
	"net"
	"log"
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

// Get room with name. Create one if no such room exists.
func (s *Server) GetRoom(n string) *Room {
	r, ok := s.Rooms[n]
	if !ok {
		r = NewRoom(n)
		s.Rooms[n] = r
	}

	return r
}

func (s *Server) DestroyIfEmpty(r *Room) {
	if r.IsEmpty() {
		delete(s.Rooms, r.Name)
	}
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

func (s *Server) RemoveUser(c *Client) error {
	for _, r := range c.Rooms {
		c.Leave(r)
		s.DestroyIfEmpty(r)
	}

	return s.DelUser(c)
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
		Prefix: s.Name,
		Command: "privmsg",
		Params: []string{text},
	}
}

func (s *Server) RoomNames() []string {
	rs := make([]string, 0, len(s.Rooms))
	for n := range s.Rooms {
		rs = append(rs, n)
	}

	return rs
}

func (s *Server) HandleClient(c *Client) {
	ws := c.Conn
	defer ws.Close()

	for {
		m := &Message{}
		if err := ws.ReadJSON(m); err != nil {
			s.RemoveUser(c)
			if isCloseError(err) {
				log.Printf("error: " + err.Error())
			}
			break
		}

		if command, ok := commands[m.Command]; ok {
			command(s, c, *m)
			continue
		}
	}
}

func isCloseError(e error) bool {
	return !websocket.IsCloseError(e,
	  websocket.CloseNormalClosure,
	  websocket.CloseGoingAway,
	  websocket.CloseNoStatusReceived)
}
