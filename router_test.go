package main

import (
	"testing"
	"reflect"
)

func TestRouterInit(t *testing.T) {
	r := new(Router).Init()

	if reflect.ValueOf(r).Elem().NumField() != 1 {
		t.Error("r does not have exactly one field")
	}

	if r.groups == nil {
		t.Error("r.groups is nil")
	}
	if len(r.groups) > 0 {
		t.Error("r.groups is not empty")
	}
}

type messageBuffer []*Message

type testClient struct {
	buf messageBuffer
}

func (c *testClient) Send(m *Message) {
	c.buf = append(c.buf, m)
}

func TestRouterManageRoutes(t *testing.T) {
	routes := map[string]clientSet{
		"/": clientSet{
			&testClient{}: struct{}{},
		},
		"/en": clientSet{
			&testClient{}: struct{}{},
		},
		"/ru": clientSet{
			&testClient{}: struct{}{},
			&testClient{}: struct{}{},
		},
	}

	r := new(Router).Init()
	for path, clients := range routes {
		for client, _ := range clients {
			r.AddRoute(path, client)
		}
	}

	if !reflect.DeepEqual(r.groups, routes) {
		t.Error("invalid routes set by r.AddRoute")
	}

	for path, clients := range routes {
		for client, _ := range clients {
			r.DelRoute(path, client)
		}
	}

	emptyRoutes := map[string]clientSet{}

	if !reflect.DeepEqual(r.groups, emptyRoutes) {
		t.Error("r.DelRoute did not destroy route entries")
	}
}

func TestRouterMessagePropogation(t *testing.T) {
	clients := map[string]*testClient{
		"sara": &testClient{
			buf: messageBuffer{},
		},
		"alex": &testClient{
			buf: messageBuffer{},
		},
		"vova": &testClient{
			buf: messageBuffer{},
		},
		"lada": &testClient{
			buf: messageBuffer{},
		},
	}

	routes := map[string][]*testClient{
		"/": []*testClient{
			clients["alex"],
		},
		"/en": []*testClient{
			clients["sara"],
		},
		"/ru": []*testClient{
			clients["vova"],
			clients["lada"],
		},
		"/private/sara": []*testClient{
			clients["sara"],
		},
	}

	messagesIndex := map[string]*Message{
		"hi /": &Message{ To: "/" },
		"hi /en": &Message{ To: "/en" },
		"hi /ru": &Message{ To: "/ru" },
		"hi /private/sara": &Message{ To: "/private/sara" },
	}

	messagesOrder := []*Message{
		messagesIndex["hi /"],
		messagesIndex["hi /en"],
		messagesIndex["hi /ru"],
		messagesIndex["hi /private/sara"],
	}

	expectedBuffers := map[string]messageBuffer{
		"alex": messageBuffer{
			messagesIndex["hi /"],
		},
		"sara": messageBuffer{
			messagesIndex["hi /en"],
			messagesIndex["hi /private/sara"],
		},
		"vova": messageBuffer{
			messagesIndex["hi /ru"],
		},
		"lada": messageBuffer{
			messagesIndex["hi /ru"],
		},
	}

	r := new(Router).Init()
	for path, clients := range routes {
		for _, client := range clients {
			r.AddRoute(path, client)
		}
	}

	for _, m := range messagesOrder {
		r.Send(m)
	}

	for name, buffer := range expectedBuffers {
		if !reflect.DeepEqual(clients[name].buf, buffer) {
			t.Errorf("%s receive wrong message sequence", name)
		}
	}
}
