package main

type Message struct {
	To string
	Body interface{}
}

type Client interface {
	Send(m *Message)
}

type clientSet map[Client]struct{}

type Router struct {
	groups map[string]clientSet
}

func (r *Router) Init() *Router {
	r.groups = map[string]clientSet{}
	return r
}

func (r *Router) Send(m *Message) {
	for c, _ := range r.groups[m.To] {
		c.Send(m)
	}
}

func (r *Router) AddRoute(p string, c Client) {
	if r.groups[p] == nil {
		r.groups[p] = clientSet{}
	}
	r.groups[p][c] = struct{}{}
}

func (r *Router) DelRoute(p string, c Client) {
	delete(r.groups[p], c)
	if len(r.groups[p]) > 0 {
		return
	}
	delete(r.groups, p)
}
