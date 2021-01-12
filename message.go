package main

import "fmt"

type Message interface {
	getMessage() string
}

type SimpleMessage struct {
	userHandle string
	msg        string
}

func (msg *SimpleMessage) getMessage() string {
	return fmt.Sprintf("%s : %s ", msg.userHandle, msg.msg)
}
