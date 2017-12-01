package main

type Message struct {
	Tags    map[string]*string
	Prefix  string
	Command string
	Params  []string
}
