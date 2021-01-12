package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("server started...")
	defer listener.Close()
	server := &Server{
		users:   make(map[string]User),
		message: make(chan Message),
	}
	go listen(server)
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			go handleConnection(server, connection)
		}

	}

}

func listen(server *Server) {
	for {
		select {
		case msg := <-server.message:
			go func(msg Message) {
				for _, value := range server.users {
					value.sendMessage(msg)
				}
			}(msg)
		}
	}
}

func handleConnection(server *Server, conn net.Conn) {
	io.WriteString(conn, "Enter your user handle : ")
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	userHandle := scanner.Text()
	user := User{
		userHandle: userHandle,
		conn:       conn,
	}

	_, ok := server.users[userHandle]
	if ok {
		io.WriteString(conn, "User handle already exists")
		conn.Close()
		return
	}
	server.users[userHandle] = user
	go waitingForMessages(server, &user)
}

func waitingForMessages(server *Server, user *User) {
	defer delete(server.users, user.userHandle) //remove user from
	defer user.conn.Close()
	sendMesageToChannel(server, user.userHandle, "Joined")
	scanner := bufio.NewScanner(user.conn)
	for scanner.Scan() {
		input := scanner.Text()
		sendMesageToChannel(server, user.userHandle, input)
	}
	sendMesageToChannel(server, user.userHandle, "Leave the chat.")
}

func sendMesageToChannel(server *Server, userHandle string, msg string) {
	server.message <- &SimpleMessage{
		userHandle: userHandle,
		msg:        msg,
	}
}
