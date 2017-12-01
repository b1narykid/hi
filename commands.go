package main

import "strings"

type Command func(*Server, *Room, *Client, []string)

var commands = map[string]Command{
	"/help":     help,
	"/who":      who,
	"/channels": channels,
	"/nick":     nick,
	"/join":     join,
	"/leave":    leave,
	"/whoami":   whoami,
	"/whois":    whois,
}

const cmds = "available commands: /help /who /channels /nick /join /leave /whoami /whois"

func wrongArgs(args []string) string {
	return "Malformatted arguments: " + strings.Join(args, " ")
}

func help(s *Server, r *Room, c *Client, args []string) {
	c.Write(s.Compose(cmds))
}

func who(s *Server, r *Room, c *Client, args []string) {
	c.Write(r.Compose(strings.Join(r.UserNames(), ", ")))
}

func channels(s *Server, r *Room, c *Client, args []string) {
	c.Write(s.Compose(strings.Join(s.RoomNames(), ", ")))
}

func nick(s *Server, r *Room, c *Client, args []string) {
	if len(args) < 2 {
		c.Conn.WriteJSON(r.Compose(wrongArgs(args)))
		return
	}

	name := args[1]
	if !UserNameIsValid(name) {
		c.Write(s.Compose(errInvalidName.Error()))
		return
	}

	oldname := c.Name
	if err := s.Nick(c.Name, name); err != nil {
		c.Write(s.Compose(err.Error()))
		return
	}

	notice := s.Compose("User " + oldname + " changed nickname to " + c.Name)
	for _, r := range c.Rooms {
		r.Write(notice)
	}
}

func join(s *Server, r *Room, c *Client, args []string) {
	if len(args) < 2 {
		c.Write(r.Compose(wrongArgs(args)))
		return
	}

	name := args[1]
	if !RoomNameIsValid(name) {
		c.Write(s.Compose(errInvalidName.Error()))
		return
	}

	room := s.GetRoom(name)
	if err := c.Join(room); err != nil {
		c.Write(room.Compose(err.Error()))
		return
	}

	greeting := room.Compose("Welcome to " + name + ", " + c.Name + "!")
	room.Write(greeting)
}

func leave(s *Server, r *Room, c *Client, args []string) {
	if len(args) < 2 {
		c.Write(r.Compose(wrongArgs(args)))
		return
	}

	name := args[1]
	room, ok := s.Rooms[name]
	if !ok {
		c.Write(s.Compose(errNoSuchRoom.Error()))
		return
	}

	if err := c.Leave(room); err != nil {
		c.Write(s.Compose("Couldn't leave " + name + ": " + err.Error()))
		return
	}

	c.Write(s.Compose("Left channel " + name + "."))
	room.Write(s.Compose("User " + c.Name + " left channel."))
}

func whoami(s *Server, r *Room, c *Client, args []string) {
	c.Write(s.Compose(c.WhoAmI()))
}

func whois(s *Server, r *Room, c *Client, args []string) {
	if len(args) < 2 {
		c.Write(r.Compose(wrongArgs(args)))
		return
	}

	name := args[1]
	if !UserNameIsValid(name) {
		c.Write(s.Compose(errInvalidName.Error()))
		return
	}

	stranger, ok := s.Users[name]
	if !ok {
		c.Write(s.Compose(errNoSuchUser.Error()))
		return
	}

	c.Write(s.Compose(stranger.WhoAmI()))
}
