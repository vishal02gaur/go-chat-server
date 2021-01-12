package main

type Server struct {
	users   map[string]User
	message chan Message
}
