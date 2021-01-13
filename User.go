package main

import (
	"io"
	"net"
)

type User struct {
	userHandle string
	conn       net.Conn
	output     chan Message
}

func (user *User) sendMessage(msg Message) {
	io.WriteString(user.conn, msg.getMessage()+"\n")
}
