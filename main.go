package main

import (
	"log"
	"flag"
	"net/http"
)

var (
	addr = ":8081"
	router = &Router{}
)

func init() {
	flag.StringVar(&addr,  "l", addr,  "bind to this addr")
	flag.Parse()
}

func main() {
	router.Init()

	http.HandleFunc("/", wsHandler)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("Server died with message: ", err)
	}
}
