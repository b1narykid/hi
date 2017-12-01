package main

type Client struct {
	Name  string
	Conn  Connection
	Rooms map[*string]*Room
}

func NewClient(name string, ws Connection) *Client {
	return &Client{
		Name: name,
		Conn: ws,
		Rooms: make(map[*string]*Room),
	}
}

func (c *Client) UniqueID() *string {
	return &c.Name
}

func (c *Client) Write(msg Message) error {
	return c.Conn.WriteJSON(msg)
}

func (c *Client) Compose(text string) Message {
	return Message{
		Prefix: c.Name,
		Command: "privmsg",
		Params: []string{text},
	}
}

func (c *Client) Leave(r *Room) error {
	rid := r.UniqueID()
	r, ok := c.Rooms[rid]
	if !ok {
		return errNoSuchUser
	}

	delete(c.Rooms, rid)
	return r.Unsubscribe(c)
}

func (c *Client) Join(r *Room) error {
	c.Rooms[r.UniqueID()] = r
	return r.Subscribe(c)
}

func (c *Client) WhoAmI() string {
	return c.Name + "@" + c.Conn.RemoteAddr().String()
}
