package main

import (
	"log"
	"flag"
	"net/http"
)

var (
	addr = ":8080"
	router = &Router{}
)

func init() {
	flag.StringVar(&addr,  "l", addr,  "bind address")
	flag.Parse()

	router.Init()
}

func main() {
	http.HandleFunc("/", wsHandler)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("Server died with message: ", err)
	}
}
