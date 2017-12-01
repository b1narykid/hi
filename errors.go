package main

import "errors"

var (
	errNeedMoreParams = errors.New("Not enough parameters")

	errRoomExists = errors.New("room already exists")
	errUserExists = errors.New("Nickname is already in use")
	errNoSuchRoom = errors.New("No such channel")
	errNoSuchUser = errors.New("No such nick/channel")
	errErroneusNickname = errors.New("Erroneus nickname")

	errNoRecipient = errors.New("No recipient given (<command>)")
	errNoTextToSend = errors.New("No text to send")

	errOpenWsConn = errors.New("could not open websocket from connection")
)
