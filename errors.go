package main

import "errors"

var (
	errRoomExists = errors.New("room already exists")
	errUserExists = errors.New("user already exists")
	errNoSuchRoom = errors.New("room does not exist")
	errNoSuchUser = errors.New("user does not exist")

	errOpenWsConn = errors.New("could not open websocket from connection")
	errInvalidName = errors.New("invalid name")
)
