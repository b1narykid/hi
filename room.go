package main

type Room struct {
	Name    string
	Clients map[*string]*Client
}

func NewRoom(name string) (r *Room) {
	return &Room{
		Name: name,
		Clients: make(map[*string]*Client),
	}
}

func (r *Room) UniqueID() *string {
	return &r.Name
}

func (c *Room) IsEmpty() bool {
	return len(c.Clients) < 1
}

func (r *Room) Write(msg Message) error {
	if r.IsEmpty() {
		return errNoSuchRoom
	}

	for _, c := range r.Clients {
		c.Write(msg)
	}

	return nil
}

func (r *Room) Compose(text string) Message {
	return Message{
		Prefix: r.Name,
		Command: "privmsg",
		Params: []string{text},
	}
}

func (r *Room) Subscribe(c *Client) error {
	cid := c.UniqueID()
	if _, ok := r.Clients[cid]; ok {
		return errUserExists
	}

	r.Clients[cid] = c
	return nil
}

func (r *Room) Unsubscribe(c *Client) error {
	cid := c.UniqueID()
	if _, ok := r.Clients[cid]; !ok {
		return errNoSuchUser
	}

	delete(r.Clients, cid)
	return nil
}

func (r *Room) UserNames() []string {
	u := make([]string, 0, len(r.Clients))
	for _, c := range r.Clients {
		u = append(u, c.Name)
	}

	return u
}
