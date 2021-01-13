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
		users:     make(map[string]User),
		broadcast: make(chan Message),
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
		case msg := <-server.broadcast:
			go func(msg Message) {
				for _, user := range server.users {
					user.output <- msg
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
		output:     make(chan Message),
	}

	_, ok := server.users[userHandle]
	if ok {
		io.WriteString(conn, "User handle already exists")
		conn.Close()
		return
	}
	server.users[userHandle] = user
	go readMessages(server, &user)
	go writeMessages(server, &user)
}

func writeMessages(server *Server, user *User) {
	for message := range user.output {
		user.sendMessage(message)
	}
}

func readMessages(server *Server, user *User) {
	defer close(user.output)
	defer user.conn.Close()
	broadcastMessage(server, user.userHandle, "Joined")
	scanner := bufio.NewScanner(user.conn)
	for scanner.Scan() {
		input := scanner.Text()
		//TODO: handle input return some garbage string. when client terminate the connection.
		broadcastMessage(server, user.userHandle, input)
	}
	delete(server.users, user.userHandle) //remove user from connected users map
	broadcastMessage(server, user.userHandle, "Leave the chat.")
}

func broadcastMessage(server *Server, userHandle string, msg string) {
	server.broadcast <- &SimpleMessage{
		userHandle: userHandle,
		msg:        msg,
	}
}
