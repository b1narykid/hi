package main

import "regexp"

var reUser = regexp.MustCompile("^(\\p{N}|\\p{L}){1,32}$")
var reRoom = regexp.MustCompile("^#(\\p{N}|\\p{L}){1,32}$")

func UserNameIsValid(n string) bool {
	return reUser.MatchString(n)
}

func RoomNameIsValid(n string) bool {
	return reRoom.MatchString(n)
}
