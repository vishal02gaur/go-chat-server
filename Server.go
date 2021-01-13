package main

type Server struct {
	users     map[string]User
	broadcast chan Message
}
